package eval

import (
	"bytes"
	"html/template"

	"github.com/gobwas/glob"
)

var _ Evaluator = new(permission)

func Permission(action string, scope string) Evaluator {
	return permission{
		Action: action,
		Scope:  scope,
	}
}

type permission struct {
	Action string
	Scope  string
}

func (p permission) Evaluate(permissions map[string]map[string]struct{}) (bool, error) {
	scopes, ok := permissions[p.Action]
	if ok {
		if p.Scope == ScopeNone {
			return true, nil
		}

		for s := range scopes {
			// TODO: replace glob parser with a simpler parser that handles only prefixes and asterisk matching.
			rule, err := glob.Compile(s, ':', '/')
			if err != nil {
				return false, err
			}
			if rule.Match(p.Scope) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (p permission) Inject(params map[string]string) (Evaluator, error) {
	tmpl, err := template.New("scope").Parse(p.Scope)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, params); err != nil {
		return nil, err
	}
	return Permission(p.Action, buf.String()), nil
}
