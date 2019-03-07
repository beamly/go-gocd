package gocd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testUnmarshallMaterialAttributesPlugin(t *testing.T) {
	m := MaterialAttributesPlugin{}
	expected := MaterialAttributesPlugin{
		Ref: "test-ref",

		Destination: "test-desintation",
		Filter: &MaterialFilter{
			Ignore: []string{"one", "two"},
		},
		InvertFilter: true,
	}

	unmarshallMaterialAttributesPlugin(&m, map[string]interface{}{
		"ref":           expected.Ref,
		"destination":   expected.Destination,
		"invert_filter": expected.InvertFilter,
		"filter": map[string]interface{}{
			"ignore": expected.Filter.Ignore,
		},
		"foo": nil,
	})

	assert.Equal(t, expected, m)
}

func testGenerateGenericPluginDependency(t *testing.T) {
	for _, tt := range []struct {
		name           string
		dependency     *MaterialAttributesPlugin
		dependencyWant map[string]interface{}
	}{
		{
			name: "basic",
			dependency: &MaterialAttributesPlugin{
				Ref:          "mock-ref",
				Destination:  "mock-destination",
				InvertFilter: true,
				Filter: &MaterialFilter{
					Ignore: []string{"one", "two"},
				},
			},
			dependencyWant: map[string]interface{}{
				"ref":           "mock-ref",
				"destination":   "mock-destination",
				"invert_filter": true,
				"filter": map[string]interface{}{
					"ignore": []interface{}{"one", "two"},
				},
			},
		},
		{
			name:           "null",
			dependency:     &MaterialAttributesPlugin{},
			dependencyWant: map[string]interface{}{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dependency.GenerateGeneric()
			assert.Equal(t, tt.dependencyWant, got)
		})
	}
}
