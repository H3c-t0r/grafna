package codegen

import (
	"bytes"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/grafana/pkg/plugins/codegen/kindsys"
)

const prefix = "github.com/grafana/grafana/public/app/plugins"

// PluginTreeListJenny creates a [codejen.ManyToOne] that produces Go code
// for loading a [pfs.TreeList] given [*kindsys.PluginDecl] as inputs.
func PluginTreeListJenny() codejen.ManyToOne[*kindsys.PluginDecl] {
	outputFile := filepath.Join("pkg", "plugins", "pfs", "corelist", "corelist_load_gen.go")

	return &ptlJenny{
		outputFile: outputFile,
		plugins:    make(map[string]bool, 0),
	}
}

type ptlJenny struct {
	outputFile string
	plugins    map[string]bool
}

func (gen *ptlJenny) JennyName() string {
	return "PluginTreeListJenny"
}

func (j *ptlJenny) Generate(decls ...*kindsys.PluginDecl) (*codejen.File, error) {
	buf := new(bytes.Buffer)
	vars := templateVars_plugin_registry{
		Header: templateVars_autogen_header{
			GenLicense: true,
		},
		Plugins: make([]struct {
			PkgName, Path, ImportPath string
			NoAlias                   bool
		}, 0, len(decls)),
	}

	type tpl struct {
		PkgName, Path, ImportPath string
		NoAlias                   bool
	}

	for _, decl := range decls {
		meta := decl.PluginMeta

		if _, exists := j.plugins[meta.Id]; exists {
			continue
		}

		vars.Plugins = append(vars.Plugins, tpl{
			PkgName:    sanitizePluginId(meta.Id),
			NoAlias:    sanitizePluginId(meta.Id) != filepath.Base(decl.PluginPath),
			ImportPath: filepath.ToSlash(filepath.Join(prefix, decl.PluginPath)),
			Path:       path.Join(append(strings.Split(prefix, "/")[3:], decl.PluginPath)...),
		})

		j.plugins[meta.Id] = true
	}

	if err := tmpls.Lookup("plugin_registry.tmpl").Execute(buf, vars); err != nil {
		return nil, fmt.Errorf("failed executing plugin registry template: %w", err)
	}

	byt, err := postprocessGoFile(genGoFile{
		path: j.outputFile,
		in:   buf.Bytes(),
	})
	if err != nil {
		return nil, fmt.Errorf("error postprocessing plugin registry: %w", err)
	}

	return codejen.NewFile(j.outputFile, byt, j), nil
}
