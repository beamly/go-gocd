package gocd

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func testServerVersionResource(t *testing.T) {
	t.Run("LessThan", testServerVersionLessThan)
	t.Run("Equal", testServerVersionEqual)
}

func testServerVersionEqual(t *testing.T) {
	for _, test := range []struct {
		v1   *ServerVersion
		v2   *ServerVersion
		want bool
	}{
		{v1: &ServerVersion{Version: "1.2.3"}, v2: &ServerVersion{Version: "1.2.3"}, want: true},
		{v1: &ServerVersion{Version: "1.2.3"}, v2: &ServerVersion{Version: "2.2.3"}, want: false},
	} {
		assert.Equal(t, test.want, test.v1.Equal(test.v2))
		assert.Equal(t, test.want, test.v2.Equal(test.v1))
	}
}

func testServerVersionLessThan(t *testing.T) {
	for _, test := range []struct {
		v1   *ServerVersion
		v2   *ServerVersion
		want bool
	}{
		{v1: &ServerVersion{Version: "1.0.0"}, v2: &ServerVersion{Version: "2.0.0"}, want: true},
		{v1: &ServerVersion{Version: "2.0.1"}, v2: &ServerVersion{Version: "2.0.0"}, want: true},
		{v1: &ServerVersion{Version: "2.0.0"}, v2: &ServerVersion{Version: "1.0.0"}, want: false},
	} {
		test.v1.parseVersion()
		test.v2.parseVersion()

		assert.Equal(t, test.want, test.v1.LessThan(test.v2))
		assert.Equal(t, !test.want, test.v2.LessThan(test.v1))
	}
}
