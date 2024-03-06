package expr

import (
	"reflect"
	"testing"

	data "github.com/grafana/grafana-plugin-sdk-go/experimental/apis/data/v0alpha1"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/schemabuilder"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/expr/classic"
	"github.com/grafana/grafana/pkg/expr/mathexp"
)

func TestQueryTypeDefinitions(t *testing.T) {
	builder, err := schemabuilder.NewSchemaBuilder(
		schemabuilder.BuilderOptions{
			PluginID: []string{DatasourceType},
			ScanCode: []schemabuilder.CodePaths{{
				BasePackage: "github.com/grafana/grafana/pkg/expr",
				CodePath:    "./",
			}},
			Enums: []reflect.Type{
				reflect.TypeOf(mathexp.ReducerSum),   // pick an example value (not the root)
				reflect.TypeOf(mathexp.UpsamplerPad), // pick an example value (not the root)
				reflect.TypeOf(ReduceModeDrop),       // pick an example value (not the root)
				reflect.TypeOf(ThresholdIsAbove),
				reflect.TypeOf(classic.ConditionOperatorAnd),
			},
		})
	require.NoError(t, err)
	err = builder.AddQueries(
		schemabuilder.QueryTypeInfo{
			Discriminators: data.NewDiscriminators("queryType", QueryTypeMath),
			GoType:         reflect.TypeOf(&MathQuery{}),
			Examples: []data.QueryExample{
				{
					Name: "constant addition",
					SaveModel: data.AsUnstructured(MathQuery{
						Expression: "$A + 10",
					}),
				},
				{
					Name: "math with two queries",
					SaveModel: data.AsUnstructured(MathQuery{
						Expression: "$A - $B",
					}),
				},
			},
		},
		schemabuilder.QueryTypeInfo{
			Discriminators: data.NewDiscriminators("queryType", QueryTypeReduce),
			GoType:         reflect.TypeOf(&ReduceQuery{}),
			Examples: []data.QueryExample{
				{
					Name: "get max value",
					SaveModel: data.AsUnstructured(ReduceQuery{
						Expression: "$A",
						Reducer:    mathexp.ReducerMax,
						Settings: &ReduceSettings{
							Mode: ReduceModeDrop,
						},
					}),
				},
			},
		},
		schemabuilder.QueryTypeInfo{
			Discriminators: data.NewDiscriminators("queryType", QueryTypeResample),
			GoType:         reflect.TypeOf(&ResampleQuery{}),
			Examples: []data.QueryExample{
				{
					Name: "resample at a every day",
					SaveModel: data.AsUnstructured(ResampleQuery{
						Expression:  "$A",
						Window:      "1d",
						Downsampler: mathexp.ReducerLast,
						Upsampler:   mathexp.UpsamplerPad,
					}),
				},
			},
		},
		schemabuilder.QueryTypeInfo{
			Discriminators: data.NewDiscriminators("queryType", QueryTypeSQL),
			GoType:         reflect.TypeOf(&SQLExpression{}),
			Examples: []data.QueryExample{
				{
					Name: "Select the first row from A",
					SaveModel: data.AsUnstructured(SQLExpression{
						Expression: "SELECT * FROM A limit 1",
					}),
				},
			},
		},
		schemabuilder.QueryTypeInfo{
			Discriminators: data.NewDiscriminators("queryType", QueryTypeClassic),
			GoType:         reflect.TypeOf(&ClassicQuery{}),
			Examples:       []data.QueryExample{
				// {
				// 	Name: "do classic query (TODO)",
				// 	SaveModel: ClassicQuery{
				// 		// ????
				// 		Conditions: []classic.ConditionJSON{},
				// 	},
				// },
			},
		},
		schemabuilder.QueryTypeInfo{
			Discriminators: data.NewDiscriminators("queryType", QueryTypeThreshold),
			GoType:         reflect.TypeOf(&ThresholdQuery{}),
			Examples:       []data.QueryExample{
				// {
				// 	Name: "TODO... a threshold query",
				// 	SaveModel: ThresholdQuery{
				// 		Expression: "$A",
				// 	},
				// },
			},
		},
	)

	require.NoError(t, err)
	_ = builder.UpdateQueryDefinition(t, "./")
}
