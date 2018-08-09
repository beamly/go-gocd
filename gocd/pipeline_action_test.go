package gocd

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func testPipelineServicePause(t *testing.T) {
	for n, test := range []struct {
		v             *ServerVersion
		confirmHeader string
		acceptHeader  string
	}{
		{
			v:             &ServerVersion{Version: "14.3.0"},
			confirmHeader: "Confirm",
			acceptHeader:  apiV0,
		},
		{
			v:             &ServerVersion{Version: "18.3.0"},
			confirmHeader: "X-GoCD-Confirm",
			acceptHeader:  apiV1,
		},
	} {
		err := test.v.parseVersion()
		assert.NoError(t, err)

		cachedServerVersion = test.v
		// defaultHTTPMux doesn't support multiple registrations so change the url a bit
		mux.HandleFunc(fmt.Sprintf("/api/pipelines/test-pipeline%d/pause", n), func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, r.Method, "POST", "Unexpected HTTP method")
			assert.Equal(t, "true", r.Header.Get(test.confirmHeader))
			if test.acceptHeader == "" {
				assert.Equal(t, len(r.Header["Accept"]), 0)
			} else {
				assert.Contains(t, r.Header["Accept"], test.acceptHeader)
			}
			fmt.Fprint(w, "")
		})
		pp, _, err := client.Pipelines.Pause(context.Background(), fmt.Sprintf("test-pipeline%d", n))
		if err != nil {
			assert.Nil(t, err)
		}
		assert.True(t, pp)
	}
}

func testPipelineServiceUnpause(t *testing.T) {
	for n, test := range []struct {
		v             *ServerVersion
		confirmHeader string
		acceptHeader  string
	}{
		{
			v:             &ServerVersion{Version: "14.3.0"},
			confirmHeader: "Confirm",
			acceptHeader:  apiV0,
		},
		{
			v:             &ServerVersion{Version: "18.3.0"},
			confirmHeader: "X-GoCD-Confirm",
			acceptHeader:  apiV1,
		},
	} {
		err := test.v.parseVersion()
		assert.NoError(t, err)

		cachedServerVersion = test.v
		// defaultHTTPMux doesn't support multiple registrations so change the url a bit
		mux.HandleFunc(fmt.Sprintf("/api/pipelines/test-pipeline%d/unpause", n), func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, r.Method, "POST", "Unexpected HTTP method")
			assert.Equal(t, "true", r.Header.Get(test.confirmHeader))
			if test.acceptHeader == "" {
				assert.Equal(t, len(r.Header["Accept"]), 0)
			} else {
				assert.Contains(t, r.Header["Accept"], test.acceptHeader)
			}
			fmt.Fprint(w, "")
		})
		pp, _, err := client.Pipelines.Unpause(context.Background(), fmt.Sprintf("test-pipeline%d", n))
		if err != nil {
			assert.Nil(t, err)
		}
		assert.True(t, pp)
	}
}

func testPipelineServiceReleaseLock(t *testing.T) {
	for n, test := range []struct {
		v             *ServerVersion
		confirmHeader string
		acceptHeader  string
	}{
		{
			v:             &ServerVersion{Version: "14.3.0"},
			confirmHeader: "Confirm",
			acceptHeader:  apiV0,
		},
		{
			v:             &ServerVersion{Version: "18.3.0"},
			confirmHeader: "X-GoCD-Confirm",
			acceptHeader:  apiV1,
		},
	} {
		err := test.v.parseVersion()
		assert.NoError(t, err)

		cachedServerVersion = test.v
		// defaultHTTPMux doesn't support multiple registrations so change the url a bit
		mux.HandleFunc(fmt.Sprintf("/api/pipelines/test-pipeline%d/releaseLock", n), func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, r.Method, "POST", "Unexpected HTTP method")
			assert.Equal(t, "true", r.Header.Get(test.confirmHeader))
			if test.acceptHeader == "" {
				assert.Equal(t, len(r.Header["Accept"]), 0)
			} else {
				assert.Contains(t, r.Header["Accept"], test.acceptHeader)
			}
			fmt.Fprint(w, "")
		})
		pp, _, err := client.Pipelines.ReleaseLock(context.Background(), fmt.Sprintf("test-pipeline%d", n))
		if err != nil {
			assert.Nil(t, err)
		}
		assert.True(t, pp)
	}
}
