package gocd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testUnmarshallMaterialAttributesPkg(t *testing.T) {
	m := MaterialAttributesPackage{}
	expected := MaterialAttributesPackage{Ref: "test-ref"}

	unmarshallMaterialAttributesPackage(&m, map[string]interface{}{"ref": expected.Ref, "foo": nil})

	assert.Equal(t, expected, m)
}

func testGenerateGenericPackageDependency(t *testing.T) {
	for _, tt := range []struct {
		name           string
		dependency     *MaterialAttributesPackage
		dependencyWant map[string]interface{}
	}{
		{
			name: "basic",
			dependency: &MaterialAttributesPackage{
				Ref: "mock-ref",
			},
			dependencyWant: map[string]interface{}{
				"ref": "mock-ref",
			},
		},
		{
			name:           "null",
			dependency:     &MaterialAttributesPackage{},
			dependencyWant: map[string]interface{}{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dependency.GenerateGeneric()
			assert.Equal(t, tt.dependencyWant, got)
		})
	}
}
