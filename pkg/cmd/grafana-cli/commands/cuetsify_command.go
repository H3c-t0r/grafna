package commands

import (
	"bytes"
	gerrors "errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/format"
	cload "cuelang.org/go/cue/load"
	"cuelang.org/go/cue/parser"
	"github.com/grafana/cuetsy"
	"github.com/grafana/grafana/pkg/cmd/grafana-cli/utils"
	"github.com/grafana/grafana/pkg/schema/load"
)

// FIXME almost this whole file is a sloppy, one-off hack that just goes around actually making
// the API we need. Parts need to be factored out appropriately.

var ctx = cuecontext.New()

const allowedImport = "github.com/grafana/grafana/packages/grafana-schema/src/schema"

var importMap = map[string]string{
	allowedImport: "@grafana-schema",
}

const prefix = "/"

func (cmd Command) generateTypescript(c utils.CommandLine) error {
	root := c.String("grafana-root")
	if root == "" {
		return gerrors.New("must provide path to the root of a Grafana repository checkout")
	}

	var fspaths load.BaseLoadPaths
	var err error

	fspaths.BaseCueFS, err = populateMapFSFromRoot(paths.BaseCueFS, root, "")
	if err != nil {
		return err
	}
	fspaths.DistPluginCueFS, err = populateMapFSFromRoot(paths.DistPluginCueFS, root, "")
	if err != nil {
		return err
	}
	overlay, err := defaultOverlay(fspaths)
	if err != nil {
		return err
	}

	// Prep the cue load config
	clcfg := &cload.Config{
		Overlay: overlay,
		// FIXME these module paths won't work for things not under our cue.mod - AKA third-party plugins
		// ModuleRoot: prefix,
		Module: "github.com/grafana/grafana",
	}

	// One-time load of the panel-plugin scuemata family def, for unifying to easily apply cuetsy attributes
	clcfg.Dir = "cue/scuemata"
	v := ctx.BuildInstance(cload.Instances(nil, clcfg)[0])
	if v.Err() != nil {
		return v.Err()
	}
	ppf := v.LookupPath(cue.ParsePath("#PanelSchema"))
	_ = ppf

	// FIXME hardcoding paths to exclude is not the way to handle this
	excl := map[string]bool{
		"cue.mod":      true,
		"cue/scuemata": true,
		"packages/grafana-schema/src/scuemata/dashboard":      true,
		"packages/grafana-schema/src/scuemata/dashboard/dist": true,
	}

	outfiles := make(map[string][]byte)

	cuetsify := func(in fs.FS) error {
		seen := make(map[string]bool)
		return fs.WalkDir(in, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			dir := filepath.Dir(path)

			if d.IsDir() || filepath.Ext(d.Name()) != ".cue" || seen[dir] || excl[dir] {
				return nil
			}
			seen[dir] = true
			clcfg.Dir = dir
			// FIXME Horrible hack to figure out the identifier used for
			// imported packages - intercept the parser called by the loader to
			// look at the ast.Files on their way in to building.
			// Much better if we could work backwards from the cue.Value,
			// maybe even directly in cuetsy itself, and figure out when a
			// referenced object is "out of bounds".
			// var imports sync.Map
			var imports []*ast.ImportSpec
			clcfg.ParseFile = func(name string, src interface{}) (*ast.File, error) {
				f, err := parser.ParseFile(name, src, parser.ParseComments)
				if err != nil {
					return nil, err
				}
				imports = append(imports, f.Imports...)
				return f, nil
			}

			// FIXME loading in this way causes all files in a dir to be loaded
			// as a single cue.Instance or cue.Value, which makes it quite
			// difficult to map them _back_ onto the original file and generate
			// discrete .gen.ts files for each .cue input.  However, going one
			// .cue file at a time and passing it as the first arg to
			// load.Instances() means that the other files are ignored
			// completely, causing references between these files to be
			// unresolved, and thus encounter a different kind of error.
			insts := cload.Instances(nil, clcfg)
			if len(insts) > 1 {
				panic("extra instances")
			}
			bi := insts[0]

			// dumpBuildInst(bi)
			v := ctx.BuildInstance(bi)
			if v.Err() != nil {
				return v.Err()
			}

			var b []byte
			f := &tsFile{}
			seen := make(map[string]bool)
			// FIXME explicitly mapping path patterns to conversion patterns
			// is exactly what we want to avoid
			switch {
			// panel plugin models.cue files
			case strings.Contains(path, "public/app/plugins"):
				for _, im := range imports {
					ip := strings.Trim(im.Path.Value, "\"")
					if ip != allowedImport {
						// TODO make a specific error type for this
						return errors.Newf(im.Pos(), "import %q not allowed, panel plugins may only import from %q", ip, allowedImport)
					}
					// TODO this approach will silently swallow the unfixable
					// error case where multiple files in the same dir import
					// the same package to a different ident
					if !seen[ip] {
						seen[ip] = true
						f.Imports = append(f.Imports, convertImport(im))
					}
				}

				// val := v.LookupPath(cue.ParsePath("Panel.lineages[0][0]"))
				// Extract the latest schema and its version number
				f.V = &tsModver{}
				lins := v.LookupPath(cue.ParsePath("Panel.lineages"))
				f.V.Lin, _ = lins.Len().Int64()
				f.V.Lin = f.V.Lin - 1
				schs := lins.LookupPath(cue.MakePath(cue.Index(int(f.V.Lin))))
				f.V.Sch, _ = schs.Len().Int64()
				f.V.Sch = f.V.Sch - 1
				latest := schs.LookupPath(cue.MakePath(cue.Index(int(f.V.Sch))))

				sch := latest.UnifyAccept(ppf, latest)
				b, err = cuetsy.Generate(sch, cuetsy.Config{})
			default:
				b, err = cuetsy.Generate(v, cuetsy.Config{})
			}

			if err != nil {
				return err
			}
			f.Body = string(b)

			var buf bytes.Buffer
			err = tsTemplate.Execute(&buf, f)
			outfiles[strings.Replace(path, ".cue", ".gen.ts", -1)] = buf.Bytes()
			return err
		})
	}

	err = cuetsify(fspaths.BaseCueFS)
	if err != nil {
		return gerrors.New(errors.Details(err, nil))
	}
	err = cuetsify(fspaths.DistPluginCueFS)
	if err != nil {
		return gerrors.New(errors.Details(err, nil))
	}

	for of, b := range outfiles {
		err := os.WriteFile(filepath.Join(root, of), b, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func convertImport(im *ast.ImportSpec) *tsImport {
	tsim := &tsImport{
		Pkg: importMap[allowedImport],
	}
	if im.Name != nil && im.Name.String() != "" {
		tsim.Ident = im.Name.String()
	} else {
		sl := strings.Split(im.Path.Value, "/")
		final := sl[len(sl)-1]
		if idx := strings.Index(final, ":"); idx != -1 {
			tsim.Pkg = final[idx:]
		} else {
			tsim.Pkg = final
		}
	}
	return tsim
}

func defaultOverlay(p load.BaseLoadPaths) (map[string]cload.Source, error) {
	overlay := make(map[string]cload.Source)

	if err := toOverlay(prefix, p.BaseCueFS, overlay); err != nil {
		return nil, err
	}

	if err := toOverlay(prefix, p.DistPluginCueFS, overlay); err != nil {
		return nil, err
	}

	return overlay, nil
}

func toOverlay(prefix string, vfs fs.FS, overlay map[string]cload.Source) error {
	if !filepath.IsAbs(prefix) {
		return fmt.Errorf("must provide absolute path prefix when generating cue overlay, got %q", prefix)
	}
	err := fs.WalkDir(vfs, ".", (func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		f, err := vfs.Open(path)
		if err != nil {
			return err
		}

		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		overlay[filepath.Join(prefix, path)] = cload.FromBytes(b)
		return nil
	}))

	if err != nil {
		return err
	}

	return nil
}

func printCUE(v cue.Value) string {
	syn := v.Syntax(
		cue.Final(),         // close structs and lists
		cue.Concrete(false), // allow incomplete values
		cue.Definitions(false),
		cue.Optional(true),
		cue.Attributes(true),
		cue.Docs(true),
	)

	bs, _ := format.Node(
		syn,
		format.TabIndent(true),
	)

	return string(bs)
}

type tsFile struct {
	V       *tsModver
	Imports []*tsImport
	Body    string
}

type tsModver struct {
	Lin, Sch int64
}

type tsImport struct {
	Ident string
	Pkg   string
}

var tsTemplate = template.Must(template.New("cuetsygen").Parse(
	`//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// This file was autogenerated by cuetsy. DO NOT EDIT!
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
{{range .Imports}}
import * as {{.Ident}} from '{{.Pkg}}';{{end}}
{{if .V}}
export const modelVersion = Object.freeze([{{ .V.Lin }}, {{ .V.Sch }}]);
{{end}}
{{.Body}}`))
