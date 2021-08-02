package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll_Evaluate(t *testing.T) {
	tests := []evaluateTestCase{
		{
			desc: "should return true for one that matches",
			evaluator: All(
				Permission("settings:write", Combine("settings", ScopeAll)),
			),
			permissions: map[string]map[string]struct{}{
				"settings:write": {"settings:**": struct{}{}},
			},
			expected: true,
		},
		{
			desc: "should return true for several that matches",
			evaluator: All(
				Permission("settings:write", Combine("settings", ScopeAll)),
				Permission("settings:read", Combine("settings", "auth.saml", ScopeAll)),
			),
			permissions: map[string]map[string]struct{}{
				"settings:write": {"settings:**": struct{}{}},
				"settings:read":  {"settings:**": struct{}{}},
			},
			expected: true,
		},
		{
			desc: "should return false if one does not match",
			evaluator: All(
				Permission("settings:write", Combine("settings", ScopeAll)),
				Permission("settings:read", Combine("settings", "auth.saml", ScopeAll)),
				Permission("report:read", Combine("reports", ScopeAll)),
			),
			permissions: map[string]map[string]struct{}{
				"settings:write": {"settings:**": struct{}{}},
				"settings:read":  {"settings:**": struct{}{}},
				"report:read":    {"report:1": struct{}{}},
			},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			ok, err := test.evaluator.Evaluate(test.permissions)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, ok)
		})
	}
}

func TestAll_Inject(t *testing.T) {
	tests := []injectTestCase{
		{
			desc:     "should inject correct param",
			expected: true,
			evaluator: All(
				Permission("reports:read", Combine("reports", Parameter(":reportId"))),
				Permission("settings:read", Combine("settings", Parameter(":settingsId"))),
			),
			params: map[string]string{
				":id":         "10",
				":settingsId": "3",
				":reportId":   "1",
			},
			permissions: map[string]map[string]struct{}{
				"reports:read": {
					"reports:1": struct{}{},
				},
				"settings:read": {
					"settings:3": struct{}{},
				},
			},
		},
		{
			desc:     "should fail for nil params",
			expected: false,
			evaluator: All(
				Permission("settings:read", Combine("reports", Parameter(":settingsId"))),
				Permission("reports:read", Combine("reports", Parameter(":reportId"))),
			),
			params: nil,
			permissions: map[string]map[string]struct{}{
				"reports:read": {
					"reports:1": struct{}{},
				},
				"settings:read": {
					"settings:3": struct{}{},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			injected, err := test.evaluator.Inject(test.params)
			assert.NoError(t, err)
			ok, err := injected.Evaluate(test.permissions)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, ok)
		})
	}
}
