package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/datasourceproxy"
	"github.com/grafana/grafana/pkg/services/datasources"
	apimodels "github.com/grafana/grafana/pkg/services/ngalert/api/tooling/definitions"
	ngmodels "github.com/grafana/grafana/pkg/services/ngalert/models"
	"github.com/grafana/grafana/pkg/web"
)

var searchRegex = regexp.MustCompile(`\{(\w+)\}`)

var NotImplementedResp = ErrResp(http.StatusNotImplemented, errors.New("endpoint not implemented"), "")

func toMacaronPath(path string) string {
	return string(searchRegex.ReplaceAllFunc([]byte(path), func(s []byte) []byte {
		m := string(s[1 : len(s)-1])
		return []byte(fmt.Sprintf(":%s", m))
	}))
}

func backendTypeByUID(ctx *models.ReqContext, cache datasources.CacheService) (apimodels.Backend, error) {
	datasourceUID := web.Params(ctx.Req)[":DatasourceUID"]
	if ds, err := cache.GetDatasourceByUID(ctx.Req.Context(), datasourceUID, ctx.SignedInUser, ctx.SkipCache); err == nil {
		switch ds.Type {
		case "loki", "prometheus":
			return apimodels.LoTexRulerBackend, nil
		case "alertmanager":
			return apimodels.AlertmanagerBackend, nil
		default:
			return 0, fmt.Errorf("%w %s", errUnexpectedBackendType, ds.Type)
		}
	}
	return 0, errBackendDoesNotExist
}

// macaron unsafely asserts the http.ResponseWriter is an http.CloseNotifier, which will panic.
// Here we impl it, which will ensure this no longer happens, but neither will we take
// advantage cancelling upstream requests when the downstream has closed.
// NB: http.CloseNotifier is a deprecated ifc from before the context pkg.
type safeMacaronWrapper struct {
	http.ResponseWriter
}

func (w *safeMacaronWrapper) CloseNotify() <-chan bool {
	return make(chan bool)
}

// replacedResponseWriter overwrites the underlying responsewriter used by a *models.ReqContext.
// It's ugly because it needs to replace a value behind a few nested pointers.
func replacedResponseWriter(ctx *models.ReqContext) (*models.ReqContext, *response.NormalResponse) {
	resp := response.CreateNormalResponse(make(http.Header), nil, 0)
	cpy := *ctx
	cpyMCtx := *cpy.Context
	cpyMCtx.Resp = web.NewResponseWriter(ctx.Req.Method, &safeMacaronWrapper{resp})
	cpy.Context = &cpyMCtx
	return &cpy, resp
}

type AlertingProxy struct {
	DataProxy *datasourceproxy.DataSourceProxyService
}

// withReq proxies a different request
func (p *AlertingProxy) withReq(
	ctx *models.ReqContext,
	method string,
	u *url.URL,
	body io.Reader,
	extractor func(*response.NormalResponse) (interface{}, error),
	headers map[string]string,
) response.Response {
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return ErrResp(http.StatusBadRequest, err, "")
	}
	for h, v := range headers {
		req.Header.Add(h, v)
	}
	newCtx, resp := replacedResponseWriter(ctx)
	newCtx.Req = req

	datasourceID := web.Params(ctx.Req)[":DatasourceID"]
	if datasourceID != "" {
		recipient, err := strconv.ParseInt(web.Params(ctx.Req)[":DatasourceID"], 10, 64)
		if err != nil {
			return ErrResp(http.StatusBadRequest, err, "DatasourceID is invalid")
		}

		p.DataProxy.ProxyDatasourceRequestWithID(newCtx, recipient)
	} else {
		datasourceUID := web.Params(ctx.Req)[":DatasourceUID"]
		if datasourceUID == "" {
			return ErrResp(http.StatusBadRequest, err, "DatasourceUID is empty")
		}
		p.DataProxy.ProxyDatasourceRequestWithUID(newCtx, datasourceUID)
	}

	status := resp.Status()
	if status >= 400 {
		errMessage := string(resp.Body())
		// if Content-Type is application/json
		// and it is successfully decoded and contains a message
		// return this as response error message
		if strings.HasPrefix(resp.Header().Get("Content-Type"), "application/json") {
			var m map[string]interface{}
			if err := json.Unmarshal(resp.Body(), &m); err == nil {
				if message, ok := m["message"]; ok {
					errMessageStr, isString := message.(string)
					if isString {
						errMessage = errMessageStr
					}
				}
			}
		} else if strings.HasPrefix(resp.Header().Get("Content-Type"), "text/html") {
			// if Content-Type is text/html
			// do not return the body
			errMessage = "redacted html"
		}
		return ErrResp(status, errors.New(errMessage), "")
	}

	t, err := extractor(resp)
	if err != nil {
		return ErrResp(http.StatusInternalServerError, err, "")
	}

	b, err := json.Marshal(t)
	if err != nil {
		return ErrResp(http.StatusInternalServerError, err, "")
	}

	return response.JSON(status, b)
}

func yamlExtractor(v interface{}) func(*response.NormalResponse) (interface{}, error) {
	return func(resp *response.NormalResponse) (interface{}, error) {
		contentType := resp.Header().Get("Content-Type")
		if !strings.Contains(contentType, "yaml") {
			return nil, fmt.Errorf("unexpected content type from upstream. expected YAML, got %v", contentType)
		}
		decoder := yaml.NewDecoder(bytes.NewReader(resp.Body()))
		decoder.KnownFields(true)

		err := decoder.Decode(v)

		return v, err
	}
}

func jsonExtractor(v interface{}) func(*response.NormalResponse) (interface{}, error) {
	if v == nil {
		// json unmarshal expects a pointer
		v = &map[string]interface{}{}
	}
	return func(resp *response.NormalResponse) (interface{}, error) {
		contentType := resp.Header().Get("Content-Type")
		if !strings.Contains(contentType, "json") {
			return nil, fmt.Errorf("unexpected content type from upstream. expected JSON, got %v", contentType)
		}
		return v, json.Unmarshal(resp.Body(), v)
	}
}

func messageExtractor(resp *response.NormalResponse) (interface{}, error) {
	return map[string]string{"message": string(resp.Body())}, nil
}

func validateCondition(ctx context.Context, c ngmodels.Condition, user *models.SignedInUser, skipCache bool, datasourceCache datasources.CacheService) error {
	if len(c.Data) == 0 {
		return nil
	}

	refIDs, err := validateQueriesAndExpressions(ctx, c.Data, user, skipCache, datasourceCache)
	if err != nil {
		return err
	}

	t := make([]string, 0, len(refIDs))
	for refID := range refIDs {
		t = append(t, refID)
	}
	if _, ok := refIDs[c.Condition]; !ok {
		return fmt.Errorf("condition %s not found in any query or expression: it should be one of: [%s]", c.Condition, strings.Join(t, ","))
	}
	return nil
}

// conditionValidator returns a curried validateCondition that accepts only condition
func conditionValidator(c *models.ReqContext, cache datasources.CacheService) func(ngmodels.Condition) error {
	return func(condition ngmodels.Condition) error {
		return validateCondition(c.Req.Context(), condition, c.SignedInUser, c.SkipCache, cache)
	}
}

func validateQueriesAndExpressions(ctx context.Context, data []ngmodels.AlertQuery, user *models.SignedInUser, skipCache bool, datasourceCache datasources.CacheService) (map[string]struct{}, error) {
	refIDs := make(map[string]struct{})
	if len(data) == 0 {
		return nil, nil
	}

	for _, query := range data {
		datasourceUID, err := query.GetDatasource()
		if err != nil {
			return nil, err
		}

		isExpression, err := query.IsExpression()
		if err != nil {
			return nil, err
		}
		if isExpression {
			refIDs[query.RefID] = struct{}{}
			continue
		}

		_, err = datasourceCache.GetDatasourceByUID(ctx, datasourceUID, user, skipCache)
		if err != nil {
			return nil, fmt.Errorf("invalid query %s: %w: %s", query.RefID, err, datasourceUID)
		}
		refIDs[query.RefID] = struct{}{}
	}
	return refIDs, nil
}

// ErrorResp creates a response with a visible error
func ErrResp(status int, err error, msg string, args ...interface{}) *response.NormalResponse {
	if msg != "" {
		err = errors.WithMessagef(err, msg, args...)
	}
	return response.Error(status, err.Error(), err)
}

// accessForbiddenResp creates a response of forbidden access.
func accessForbiddenResp() response.Response {
	return ErrResp(http.StatusForbidden, errors.New("Permission denied"), "")
}
