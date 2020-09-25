package eval_test

import (
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"testing"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/internal/eval"
	"github.com/b2wdigital/restQL-golang/v4/test"
)

func TestApplyAggregators(t *testing.T) {
	tests := []struct {
		name      string
		query     domain.Query
		resources domain.Resources
		expected  domain.Resources
	}{
		{
			"should do nothing if there is no aggregator",
			domain.Query{Statements: []domain.Statement{{Resource: "hero"}, {Resource: "sidekick"}}},
			domain.Resources{
				"hero":     restql.DoneResource{ResponseBody: nil},
				"sidekick": restql.DoneResource{ResponseBody: nil},
			},
			domain.Resources{
				"hero":     restql.DoneResource{ResponseBody: nil},
				"sidekick": restql.DoneResource{ResponseBody: nil},
			},
		},
		{
			"should aggregate one resource inside other",
			domain.Query{Statements: []domain.Statement{
				{Resource: "hero"},
				{Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
			domain.Resources{
				"hero":     restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`{ "id": 1, "name": "batman" }`),
				)},
				"sidekick": restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`{ "id": 10, "name": "robin" }`),
				)},
			},
			domain.Resources{
				"hero":     restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`{ "id": 1, "name": "batman", "sidekick": { "id": 10, "name": "robin" } }`),
				)},
				"sidekick": restql.DoneResource{ResponseBody: &restql.ResponseBody{}},
			},
		},
		{
			"should aggregate one resource inside other in deep location",
			domain.Query{Statements: []domain.Statement{
				{Resource: "hero"},
				{Resource: "sidekick", In: []string{"hero", "info", "partners", "sidekick"}},
			}},
			domain.Resources{
				"hero":     restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`{ "id": 1, "name": "batman" }`),
				)},
				"sidekick": restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`{ "id": 10, "name": "robin" }`),
				)},
			},
			domain.Resources{
				"hero":     restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`{ "id": 1, "name": "batman", "info": { "partners": { "sidekick": { "id": 10, "name": "robin" } } } }`),
				)},
				"sidekick": restql.DoneResource{ResponseBody: &restql.ResponseBody{}},
			},
		},
		{
			"should aggregate one list resource inside other",
			domain.Query{Statements: []domain.Statement{
				{Resource: "hero"},
				{Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
			domain.Resources{
				"hero":     restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`{ "id": 1, "name": "batman" }`),
				)},
				"sidekick": restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`[{ "id": 10, "name": "robin" }, { "id": 11, "name": "batgirl" }]`),
				)},
			},
			domain.Resources{
				"hero":     restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`{ "id": 1, "name": "batman", "sidekick":  [{ "id": 10, "name": "robin" }, { "id": 11, "name": "batgirl" }]}`),
				)},
				"sidekick": restql.DoneResource{ResponseBody: &restql.ResponseBody{}},
			},
		},
		{
			"should aggregate one resource inside other multiplexed resource",
			domain.Query{Statements: []domain.Statement{
				{Resource: "hero"},
				{Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
			domain.Resources{
				"hero": restql.DoneResources{
					restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": 1, "name": "batman" }`),
					)},
					restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": 2, "name": "wonder woman" }`),
					)},
				},
				"sidekick": restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`{ "id": 10, "name": "robin" }`),
				)},
			},
			domain.Resources{
				"hero": restql.DoneResources{
					restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": 1, "name": "batman", "sidekick": { "id": 10, "name": "robin" } }`),
					)},
					restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": 2, "name": "wonder woman", "sidekick": { "id": 10, "name": "robin" } }`),
					)},
				},
				"sidekick": restql.DoneResource{ResponseBody: &restql.ResponseBody{}},
			},
		},
		{
			"should aggregate one resource inside every item on target result",
			domain.Query{Statements: []domain.Statement{
				{Resource: "hero"},
				{Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
			domain.Resources{
				"hero":     restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`[{ "id": 1, "name": "batman" }, { "id": 2, "name": "wonder woman" }]`),
				)},
				"sidekick": restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`{ "id": 10, "name": "robin" }`),
				)},
			},
			domain.Resources{
				"hero":     restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`[{ "id": 1, "name": "batman", "sidekick": { "id": 10, "name": "robin" } }, { "id": 2, "name": "wonder woman", "sidekick": { "id": 10, "name": "robin" } }]`),
				)},
				"sidekick": restql.DoneResource{ResponseBody: &restql.ResponseBody{}},
			},
		},
		{
			"should aggregate one multiplexed resource inside other",
			domain.Query{Statements: []domain.Statement{
				{Resource: "hero"},
				{Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
			domain.Resources{
				"hero": restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`{ "id": 1, "name": "batman" }`),
				)},
				"sidekick": restql.DoneResources{
					restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": 10, "name": "robin" }`),
					)},
					restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": 11, "name": "batgirl" }`),
					)},
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": 1, "name": "batman", "sidekick": [{ "id": 10, "name": "robin" }, { "id": 11, "name": "batgirl" }] }`),
					),
				},
				"sidekick": restql.DoneResources{
					restql.DoneResource{ResponseBody: &restql.ResponseBody{}},
					restql.DoneResource{ResponseBody: &restql.ResponseBody{}},
				},
			},
		},
		{
			"should aggregate one multiplexed resource with other multiplexed resource",
			domain.Query{Statements: []domain.Statement{
				{Resource: "hero"},
				{Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
			domain.Resources{
				"hero": restql.DoneResources{
					restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": 1, "name": "batman"}`),
					)},
					restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": 2, "name": "wonder woman"}`),
					)},
				},
				"sidekick": restql.DoneResources{
					restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": 10, "name": "robin" }`),
					)},
					restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": 11, "name": "batgirl" }`),
					)},
				},
			},
			domain.Resources{
				"hero": restql.DoneResources{
					restql.DoneResource{
						ResponseBody: restql.NewResponseBodyFromValue(
							test.NoOpLogger,
							test.Unmarshal(`{ "id": 1, "name": "batman", "sidekick": { "id": 10, "name": "robin" } }`),
						),
					},
					restql.DoneResource{
						ResponseBody: restql.NewResponseBodyFromValue(
							test.NoOpLogger,
							test.Unmarshal(`{ "id": 2, "name": "wonder woman", "sidekick": { "id": 11, "name": "batgirl" } }`),
						),
					},
				},
				"sidekick": restql.DoneResources{
					restql.DoneResource{ResponseBody: &restql.ResponseBody{}},
					restql.DoneResource{ResponseBody: &restql.ResponseBody{}},
				},
			},
		},
		{
			"should aggregate one list resource with other list resource",
			domain.Query{Statements: []domain.Statement{
				{Resource: "hero"},
				{Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
			domain.Resources{
				"hero":     restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`[{ "id": 1, "name": "batman" }, { "id": 2, "name": "wonder woman" }]`),
				)},
				"sidekick": restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`[{ "id": 10, "name": "robin" }, { "id": 11, "name": "batgirl" }]`),
				)},
			},
			domain.Resources{
				"hero":     restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(
					test.NoOpLogger,
					test.Unmarshal(`[{ "id": 1, "name": "batman", "sidekick": { "id": 10, "name": "robin" } }, { "id": 2, "name": "wonder woman", "sidekick": { "id": 11, "name": "batgirl" } }]`),
				)},
				"sidekick": restql.DoneResource{ResponseBody: &restql.ResponseBody{}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := eval.ApplyAggregators(tt.query, tt.resources)
			test.Equal(t, got, tt.expected)
		})
	}
}
