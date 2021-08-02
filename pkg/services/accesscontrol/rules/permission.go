package rules

import (
	"bytes"
	"html/template"

	"github.com/gobwas/glob"
)

var _ Evaluator = new(permission)

func Permission(action string, scopes ...string) Evaluator {
	return permission{
		Action: action,
		Scopes: scopes,
	}
}

type permission struct {
	Action string
	Scopes []string
}

func (p permission) Evaluate(permissions map[string]map[string]struct{}) (bool, error) {
	targetScopes, ok := permissions[p.Action]
	if !ok {
		return false, nil
	}

	if len(p.Scopes) == 0 {
		return true, nil
	}

	for _, scope := range p.Scopes {
		var err error
		var matches bool

		for target := range targetScopes {
			matches, err = match(scope, target)
			if err != nil {
				return false, err
			}
			if matches {
				break
			}
		}
		if !matches {
			return false, nil
		}
	}

	return true, nil
}

func match(scope, target string) (bool, error) {
	// TODO: replace glob parser with a simpler parser that handles only prefixes and asterisk matching.
	rule, err := glob.Compile(target, ':', '/')
	if err != nil {
		return false, err
	}
	if rule.Match(scope) {
		return true, nil
	}
	return false, nil
}

func (p permission) Inject(params map[string]string) (Evaluator, error) {
	scopes := make([]string, 0, len(p.Scopes))
	for _, scope := range p.Scopes {
		tmpl, err := template.New("scope").Parse(scope)
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		if err = tmpl.Execute(&buf, params); err != nil {
			return nil, err
		}
		scopes = append(scopes, buf.String())
	}
	return Permission(p.Action, scopes...), nil
}
