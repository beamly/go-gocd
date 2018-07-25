package gocd

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/hashicorp/go-version"
)

func testServerVersionResource(t *testing.T) {
	t.Run("LessThan", testServerVersionLessThan)
	t.Run("Equal", testServerVersionEqual)
	t.Run("GetAPIVersion", testServerVersionGetAPIVersion)
	t.Run("GetAPIVersionFail", testServerVersionGetAPIVersionFail)
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
		{v1: &ServerVersion{Version: "2.0.1"}, v2: &ServerVersion{Version: "2.0.0"}, want: false},
		{v1: &ServerVersion{Version: "2.0.0"}, v2: &ServerVersion{Version: "2.0.1"}, want: true},
		{v1: &ServerVersion{Version: "2.0.0"}, v2: &ServerVersion{Version: "1.0.0"}, want: false},
	} {
		name := fmt.Sprintf("%s < %s = %t", test.v1.Version, test.v2.Version, test.want)
		t.Run(name, func(t *testing.T) {

			test.v1.parseVersion()
			test.v2.parseVersion()

			assert.Equal(t, test.want, test.v1.LessThan(test.v2))
			assert.Equal(t, !test.want, test.v2.LessThan(test.v1))
		})
	}
}

func testServerVersionGetAPIVersion(t *testing.T) {
	for _, test := range []struct {
		v        *ServerVersion
		endpoint string
		method   string
		want     string
	}{
		{
			endpoint: "/api/version",
			method:   http.MethodGet,
			want:     apiV1,
			v:        &ServerVersion{Version: "1.0.0"},
		},
	} {
		test.v.parseVersion()
		apiV, err := test.v.GetAPIVersion(test.endpoint, test.method)

		assert.NoError(t, err)
		assert.Equal(t, apiV, test.want)
	}
}

func testServerVersionGetAPIVersionFail(t *testing.T) {
	for _, test := range []struct {
		v        *ServerVersion
		endpoint string
		method   string
		want     string
	}{
		{
			endpoint: "/api/foobar",
			method:   http.MethodGet,
			want:     "could not find API version tag for 'GET /api/foobar'",
		},
	} {
		apiV, err := test.v.GetAPIVersion(test.endpoint, test.method)

		assert.EqualError(t, err, test.want)
		assert.Empty(t, apiV)
	}
}

func TestNewServerApiVersionMapping(t *testing.T) {

	mockVersion, err := version.NewVersion("1.0.0")
	assert.NoError(t, err)
	type args struct {
		serverVersion string
		apiVersion    string
	}
	tests := []struct {
		name        string
		args        args
		wantMapping *ServerApiVersionMapping
	}{
		{
			name: "base",
			args: args{serverVersion: "1.0.0", apiVersion: apiV1},
			wantMapping: &ServerApiVersionMapping{
				Api:    apiV1,
				Server: mockVersion,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t,
				tt.wantMapping,
				NewServerApiVersionMapping(tt.args.serverVersion, tt.args.apiVersion),
			)
		})
	}
}
