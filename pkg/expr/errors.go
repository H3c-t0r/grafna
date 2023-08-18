package expr

import (
	"errors"
	"fmt"

	"github.com/grafana/grafana/pkg/util/errutil"
)

var ConversionError = errutil.NewBase(
	errutil.StatusBadRequest,
	"sse.readDataError",
).MustTemplate(
	"[{{ .Public.refId }}] got error: {{ .Error }}",
	errutil.WithPublic(
		"failed to read data from from query {{ .Public.refId }}: {{ .Public.error }}",
	),
)

func MakeConversionError(refID string, err error) error {
	data := errutil.TemplateData{
		// Conversion errors should only have meta information in errors
		Public: map[string]interface{}{
			"refId": refID,
			"error": err.Error(),
		},
		Error: err,
	}
	return ConversionError.Build(data)
}

var QueryError = errutil.NewBase(
	errutil.StatusBadRequest, "sse.dataQueryError").MustTemplate(
	"failed to execute query [{{ .Public.refId }}]: {{ .Error }}",
	errutil.WithPublic(
		"failed to execute query [{{ .Public.refId }}]: {{ .Public.error }}",
	))

func MakeQueryError(refID, datasourceUID string, err error) error {
	var pErr error
	var utilErr errutil.Error
	// See if this is grafana error, if so, grab public message
	if errors.As(err, &utilErr) {
		pErr = utilErr.Public()
	} else {
		pErr = err
	}

	data := errutil.TemplateData{
		Public: map[string]interface{}{
			"refId":         refID,
			"datasourceUID": datasourceUID,
			"error":         pErr.Error(),
		},
		Error: err,
	}

	return QueryError.Build(data)
}

var depErrStr = "did not execute expression [{{ .Public.refId }}] due to a failure to of the dependent expression or query [{{.Public.depRefId}}]"

var DependencyError = errutil.NewBase(
	errutil.StatusBadRequest, "sse.dependencyError").MustTemplate(
	depErrStr,
	errutil.WithPublic(
		depErrStr,
	))

func MakeDependencyError(refID, depRefID string) error {
	data := errutil.TemplateData{
		Public: map[string]interface{}{
			"refId":    refID,
			"depRefId": depRefID,
		},
		Error: fmt.Errorf("did not execute expression %v due to a failure to of the dependent expression or query %v", refID, depRefID),
	}

	return DependencyError.Build(data)
}
