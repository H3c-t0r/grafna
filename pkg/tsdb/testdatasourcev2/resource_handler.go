package testdatasourcev2

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/grafana/grafana/pkg/infra/log"

	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
)

func (td *testdataV2) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", td.testGetHandler)
	mux.HandleFunc("/stream", td.testStreamHandler)
	mux.Handle("/test", createJSONHandler(td.logger))
	mux.Handle("/test/json", createJSONHandler(td.logger))
	mux.HandleFunc("/boom", td.testPanicHandler)
}

func (td *testdataV2) testGetHandler(rw http.ResponseWriter, req *http.Request) {
	td.logger.Debug("Received resource call", "url", req.URL.String(), "method", req.Method)

	if req.Method != http.MethodGet {
		return
	}

	if _, err := rw.Write([]byte("Hello world from test datasource!")); err != nil {
		td.logger.Error("Failed to write response", "error", err)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (td *testdataV2) testStreamHandler(rw http.ResponseWriter, req *http.Request) {
	td.logger.Debug("Received resource call", "url", req.URL.String(), "method", req.Method)

	if req.Method != http.MethodGet {
		return
	}

	count := 10
	countstr := req.URL.Query().Get("count")
	if countstr != "" {
		if i, err := strconv.Atoi(countstr); err == nil {
			count = i
		}
	}

	sleep := req.URL.Query().Get("sleep")
	sleepDuration, err := time.ParseDuration(sleep)
	if err != nil {
		sleepDuration = time.Millisecond
	}

	rw.Header().Add("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)

	for i := 1; i <= count; i++ {
		if _, err := io.WriteString(rw, fmt.Sprintf("Message #%d", i)); err != nil {
			td.logger.Error("Failed to write response", "error", err)
			return
		}
		rw.(http.Flusher).Flush()
		time.Sleep(sleepDuration)
	}
}

func createJSONHandler(logger log.Logger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		logger.Debug("Received resource call", "url", req.URL.String(), "method", req.Method)

		var reqData map[string]interface{}
		if req.Body != nil {
			defer req.Body.Close()
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				logger.Error("Failed to read request body to bytes", "error", err)
			} else {
				err := json.Unmarshal(b, &reqData)
				if err != nil {
					logger.Error("Failed to unmarshal request body to JSON", "error", err)
				}

				logger.Debug("Received resource call body", "body", reqData)
			}
		}

		config := httpadapter.PluginConfigFromContext(req.Context())

		data := map[string]interface{}{
			"message": "Hello world from test datasource!",
			"request": map[string]interface{}{
				"method":  req.Method,
				"url":     req.URL,
				"headers": req.Header,
				"body":    reqData,
				"config":  config,
			},
		}
		bytes, err := json.Marshal(&data)
		if err != nil {
			logger.Error("Failed to marshal response body to JSON", "error", err)
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		if _, err := rw.Write(bytes); err != nil {
			logger.Error("Failed to write response", "error", err)
		}
	})
}

func (td *testdataV2) testPanicHandler(rw http.ResponseWriter, req *http.Request) {
	panic("BOOM")
}
