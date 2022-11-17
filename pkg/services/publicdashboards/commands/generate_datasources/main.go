//go:build ignore

package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/grafana/grafana/pkg/services/publicdashboards/commands/generate_datasources"
)

var tsDatasourcesTemplate = `
// Code generated by go generate; DO NOT EDIT.

export const supportedDatasources: { [key: string]: boolean } = {
{{- range . }}
	"{{ printf "%s" . }}": true,
{{- end }}
}
`

func main() {
	baseUrl := "https://grafana.com"
	slugs, err := generate_datasources.GetCompatibleDatasources(baseUrl)
	tsTemplate := template.Must(template.New("").Parse(tsDatasourcesTemplate))

	// Generate supported datasources for Typescript
	tsFile, err := os.Create("./../../../public/app/features/dashboard/components/ShareModal/SharePublicDashboard/SupportedPubdashDatasources.ts")
	if err != nil {
		fmt.Println(err)
	}
	err = tsTemplate.Execute(tsFile, slugs)
	if err != nil {
		fmt.Println(err)
	}
}
