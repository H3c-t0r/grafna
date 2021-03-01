package classic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana/pkg/expr/mathexp"
)

// ConditionsCmd is command for the classic conditions
// expression operation.
type ConditionsCmd struct {
	Conditions []condition
	refID      string
}

// classicConditionJSON is the JSON model for a single condition.
// It is based on services/alerting/conditions/query.go's newQueryCondition().
type classicConditionJSON struct {
	Evaluator conditionEvalJSON `json:"evaluator"`

	Operator struct {
		Type string `json:"type"`
	} `json:"operator"`

	Query struct {
		Params []string
	} `json:"query"`

	Reducer struct {
		// Params []interface{} `json:"params"` (Unused)
		Type string `json:"type"`
	}
}

type conditionEvalJSON struct {
	Params []float64 `json:"params"`
	Type   string    `json:"type"` // e.g. "gt"

}

// condition is a single condition within the ConditionsCmd.
type condition struct {
	QueryRefID string
	Reducer    classicReducer
	Evaluator  evaluator
	Operator   string
}

type classicReducer string

// NeedsVars returns the variable names (refIds) that are dependencies
// to execute the command and allows the command to fulfill the Command interface.
func (ccc *ConditionsCmd) NeedsVars() []string {
	vars := []string{}
	for _, c := range ccc.Conditions {
		vars = append(vars, c.QueryRefID)
	}
	return vars
}

// Execute runs the command and returns the results or an error if the command
// failed to execute.
func (ccc *ConditionsCmd) Execute(ctx context.Context, vars mathexp.Vars) (mathexp.Results, error) {
	firing := true
	newRes := mathexp.Results{}

	for i, c := range ccc.Conditions {
		querySeriesSet := vars[c.QueryRefID]
		for _, val := range querySeriesSet.Values {
			series, ok := val.(mathexp.Series)
			if !ok {
				return newRes, fmt.Errorf("can only reduce type series, got type %v", val.Type())
			}

			reducedNum := c.Reducer.Reduce(series)
			// TODO handle error / no data signals
			evalRes := c.Evaluator.Eval(reducedNum)

			if i == 0 {
				firing = evalRes
				//noDataFound = cr.NoDataFound
			}

			if c.Operator == "or" {
				firing = firing || evalRes
			} else {
				firing = firing && evalRes
			}
		}
	}

	num := mathexp.NewNumber("", nil)

	var v float64
	if firing {
		v = 1
	}

	num.SetValue(&v)

	newRes.Values = append(newRes.Values, num)

	return newRes, nil
}

// UnmarshalConditionsCmd creates a new ConditionsCmd.
func UnmarshalConditionsCmd(rawQuery map[string]interface{}, refID string) (*ConditionsCmd, error) {
	jsonFromM, err := json.Marshal(rawQuery["conditions"])
	if err != nil {
		return nil, fmt.Errorf("failed to remarshal classic condition body: %w", err)
	}
	var ccj []classicConditionJSON
	if err = json.Unmarshal(jsonFromM, &ccj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal remarshaled classic condition body: %w", err)
	}

	c := &ConditionsCmd{
		refID: refID,
	}

	for i, cj := range ccj {
		cond := condition{}

		if len(cj.Query.Params) == 0 || cj.Query.Params[0] == "" {
			return nil, fmt.Errorf("classic condition %v is missing the query refID argument", i+1)
		}

		cond.QueryRefID = cj.Query.Params[0]

		cond.Reducer = classicReducer(cj.Reducer.Type)
		if !cond.Reducer.ValidReduceFunc() {
			return nil, fmt.Errorf("reducer '%v' in condition %v is not a valid reducer", cond.Reducer, i+1)
		}

		cond.Evaluator, err = newAlertEvaluator(cj.Evaluator)
		if err != nil {
			return nil, err
		}

		c.Conditions = append(c.Conditions, cond)
	}

	return c, nil
}
