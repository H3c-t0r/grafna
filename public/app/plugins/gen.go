//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/grafana/grafana/pkg/codegen"
	"github.com/grafana/grafana/pkg/cuectx"
	"github.com/grafana/grafana/pkg/plugins/pfs"
)

var skipPlugins = map[string]bool{
	"barchart":       true,
	"canvas":         true,
	"histogram":      true,
	"heatmap":        true,
	"heatmap-old":    true,
	"candlestick":    true,
	"state-timelin":  true,
	"status-history": true,
	"table":          true,
	"timeseries":     true,
}

// Generate TypeScript for all plugin models.cue
func main() {
	if len(os.Args) > 1 {
		fmt.Fprintf(os.Stderr, "plugin thema code generator does not currently accept any arguments\n, got %q", os.Args)
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get working directory: %s", err)
		os.Exit(1)
	}
	//
	// var find func(path string) (string, error)
	// find = func(path string) (string, error) {
	// 	parent := filepath.Dir(path)
	// 	if parent == path {
	// 		return "", errors.New("grafana root directory could not be found")
	// 	}
	// 	fp := filepath.Join(path, "go.mod")
	// 	if _, err := os.Stat(fp); err == nil {
	// 		return path, nil
	// 	}
	// 	return find(parent)
	// }
	// groot, err := find(cwd)
	// if err != nil {
	// 	fmt.Fprint(os.Stderr, err)
	// 	os.Exit(1)
	// }

	wd := codegen.NewWriteDiffer()
	parent := os.DirFS(cwd)
	lib := cuectx.ProvideThemaLibrary()
	ptrees := make(map[string]*pfs.Tree)
	for _, typ := range []string{"datasource", "panel"} {
		ents, err := fs.ReadDir(parent, typ)
		if err != nil {
			fmt.Fprintf(os.Stderr, "readdir failed on %s: %w", filepath.Join(cwd, typ), err)
			os.Exit(1)
		}
		for _, plugdir := range ents {
			if skipPlugins[plugdir.Name()] {
				continue
			}
			subp := filepath.Join(typ, plugdir.Name())
			fullp := filepath.Join(cwd, subp)
			sub, err := fs.Sub(parent, subp)
			if err != nil {
				fmt.Fprintf(os.Stderr, "creating subfs failed on %s: %w", fullp, err)
				os.Exit(1)
			}
			ptree, err := pfs.ParsePluginFS(sub, lib)
			if err != nil {
				fmt.Fprintf(os.Stderr, "creating subfs failed on %s: %w", fullp, err)
				os.Exit(1)
			}
			ptrees[fullp] = ptree
		}
	}

	for fullp, ptree := range ptrees {
		twd, err := codegen.CuetsifyPlugin(ptree, fullp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "generating typescript failed for %s: %s", fullp, err)
			os.Exit(1)
		}
		wd.Merge(twd)
	}

	// wd, err := codegen.CuetsifyPlugins(cuecontext.New(), groot)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "error while generating code:\n%s\n", err)
	// 	os.Exit(1)
	// }

	if _, set := os.LookupEnv("CODEGEN_VERIFY"); set {
		err = wd.Verify()
		if err != nil {
			fmt.Fprintf(os.Stderr, "generated code is out of sync with inputs:\n%s\nrun `make gen-cue` to regenerate\n\n", err)
			os.Exit(1)
		}
	} else {
		err = wd.Write()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while writing generated code to disk:\n%s\n", err)
			os.Exit(1)
		}
	}
}
