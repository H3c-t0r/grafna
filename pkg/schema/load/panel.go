package load

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"reflect"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	"github.com/grafana/grafana/pkg/schema"
)

// TODO a proper approach to this type would probably include responsibility for
// moving between the declared and composed forms of panel plugin schema
type panelSchema struct {
	actual    cue.Value
	major     int
	minor     int
	next      *panelSchema
	migration migrationFunc
}

func (ps *panelSchema) Validate(r schema.Resource) error {
	rv, err := rt.Compile("resource", r.Value)
	if err != nil {
		return err
	}
	return ps.actual.Unify(rv.Value()).Validate(cue.Concrete(true))
}

func (ps *panelSchema) ApplyDefaults(_ schema.Resource) (schema.Resource, error) {
	panic("not implemented") // TODO: Implement
}

func (ps *panelSchema) TrimDefaults(_ schema.Resource) (schema.Resource, error) {
	panic("not implemented") // TODO: Implement
}

func (ps *panelSchema) CUE() cue.Value {
	return ps.actual
}

func (ps *panelSchema) Version() (major int, minor int) {
	return ps.major, ps.minor
}

func (ps *panelSchema) Successor() (schema.VersionedCueSchema, bool) {
	panic("not implemented")
}

func (ps *panelSchema) Migrate(x schema.Resource) (schema.Resource, schema.VersionedCueSchema, error) {
	panic("not implemented")
}

// Returns a disjunction of structs representing each panel schema version
// (post-mapping from on-disk #PanelModel form) from each scuemata in the map.
func disjunctPanelScuemata(scuemap map[string]schema.Fam) (cue.Value, error) {
	partsi, err := rt.Compile("panelDisjunction", `
	allPanels: [Name=_]: {}
	parts: or([for v in allPanels { v }])
	`)
	if err != nil {
		return cue.Value{}, err
	}

	parts := partsi.Value()
	for id, fam := range scuemap {
		sch := fam.First()

		// TODO lol, be better
		for !reflect.ValueOf(sch).IsNil() {
			cv := mapPanelModel(id, sch)

			mjv, miv := sch.Version()
			parts = parts.Fill(cv, "allPanels", fmt.Sprintf("%s@%v.%v", id, mjv, miv))
			sch = sch.Successor()
		}
	}

	return parts.LookupPath(cue.MakePath(cue.Str("parts"))), nil
}

// mapPanelModel maps a schema from the #PanelModel form in which it's declared
// in a plugin's model.cue to the structure in which it actually appears in the
// dashboard schema.
func mapPanelModel(id string, vcs schema.VersionedCueSchema) cue.Value {
	maj, min := vcs.Version()
	// Ignore err return, this can't fail to compile
	inter, _ := rt.Compile("typedPanel", fmt.Sprintf(`
	in: {
		type: %q
		v: {
			maj: %d
			min: %d
		}
		model: {...}
	}
	result: {
		type: in.type,
		panelSchema: maj: in.v.maj
		panelSchema: min: in.v.min
		options: in.model.PanelOptions
		fieldConfig: defaults: custom: in.model.PanelFieldConfig
	}
	`, id, maj, min))

	// TODO validate, especially with #PanelModel
	return inter.Value().Fill(vcs.CUE(), "in", "model").Lookup("result")
}

func readPanelModels(p BaseLoadPaths) (map[string]schema.Fam, error) {
	overlay := make(map[string]load.Source)
	if err := toOverlay("/", p.BaseCueFS, overlay); err != nil {
		return nil, err
	}
	if err := toOverlay("/", p.DistPluginCueFS, overlay); err != nil {
		return nil, err
	}

	cfg := &load.Config{
		Overlay: overlay,
		Package: "scuemata",
	}

	built, err := rt.Build(load.Instances([]string{"/cue/scuemata/scuemata.cue"}, cfg)[0])
	if err != nil {
		return nil, err
	}
	pmf := built.Value().LookupPath(cue.MakePath(cue.Def("#PanelModelFamily")))
	if !pmf.Exists() {
		return nil, errors.New("could not locate #PanelModelFamily definition")
	}

	all := make(map[string]schema.Fam)
	err = fs.WalkDir(p.DistPluginCueFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || d.Name() != "plugin.json" {
			return nil
		}

		dpath := filepath.Dir(path)
		// For now, skip plugins without a models.cue
		_, err = p.DistPluginCueFS.Open(filepath.Join(dpath, "models.cue"))
		if err != nil {
			return nil
		}

		fi, err := p.DistPluginCueFS.Open(path)
		if err != nil {
			return err
		}
		b, err := ioutil.ReadAll(fi)
		if err != nil {
			return err
		}

		jmap := make(map[string]interface{})
		json.Unmarshal(b, &jmap)
		if err != nil {
			return err
		}
		iid, has := jmap["id"]
		if !has || jmap["type"] != "panel" {
			return errors.New("no type field in plugin.json or not a panel type plugin")
		}
		id := iid.(string)

		cfg := &load.Config{
			Package: "grafanaschema",
			Overlay: overlay,
		}

		li := load.Instances([]string{filepath.Join("/", dpath, "models.cue")}, cfg)
		built := cue.Build(li)
		// TODO this is a silly check...right?
		if len(built) != 1 {
			return fmt.Errorf("expected exactly one instance, got %v", len(built))
		}
		imod := built[0]

		// Verify that there exists a Model declaration in the models.cue file...
		// TODO Best (?) ergonomics for entire models.cue file to emit a struct
		// compliant with #PanelModelFamily
		pmod := imod.Lookup("Model")
		if !pmod.Exists() {
			return fmt.Errorf("%s does not contain a declaration of its models at path 'Model'", path)
		}

		// Ensure the declared value is subsumed by/correct wrt #PanelModelFamily
		if err := pmf.Subsume(pmod); err != nil {
			return err
		}

		// Create a generic schema family to represent the whole of the
		fam, err := buildGenericScuemata(pmod)
		if err != nil {
			return err
		}

		all[id] = fam
		return nil
	})
	if err != nil {
		return nil, err
	}

	return all, nil
}
