package gocd

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestServerVersion(t *testing.T) {
	t.Run("ServerVersionCaching", testServerVersionCaching)
	t.Run("ServerVersion", testServerVersion)
	t.Run("Resource", testServerVersionResource)
}

func testServerVersion(t *testing.T) {
	setup()
	defer teardown()
	ver, err := version.NewVersion("16.6.0")
	assert.NoError(t, err)

	mux.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		assert.Equal(t, apiV1, r.Header.Get("Accept"))

		j, _ := ioutil.ReadFile("test/resources/server-version.v1.1.json")

		fmt.Fprint(w, string(j))
	})

	cachedServerVersion = nil
	v, _, err := client.ServerVersion.Get(context.Background())

	assert.NoError(t, err)

	assert.Equal(t, &ServerVersion{
		Version:      "16.6.0",
		BuildNumber:  "3348",
		GitSha:       "a7a5717cbd60c30006314fb8dd529796c93adaf0",
		FullVersion:  "16.6.0 (3348-a7a5717cbd60c30006314fb8dd529796c93adaf0)",
		CommitURL:    "https://github.com/gocd/gocd/commits/a7a5717cbd60c30006314fb8dd529796c93adaf0",
		VersionParts: ver,
	}, v)

	// Verify that the server version is cached
	assert.Equal(t, cachedServerVersion, v)

}

func testServerVersionCaching(t *testing.T) {
	if runIntegrationTest() {
		ver, err := version.NewVersion("18.7.0")
		assert.NoError(t, err)

		cachedServerVersion = &ServerVersion{
			Version:      "18.7.0",
			BuildNumber:  "7121",
			GitSha:       "75d1247f58ab8bcde3c5b43392a87347979f82c5",
			FullVersion:  "18.7.0 (7121-75d1247f58ab8bcde3c5b43392a87347979f82c5)",
			CommitURL:    "https://github.com/gocd/gocd/commits/75d1247f58ab8bcde3c5b43392a87347979f82c5",
			VersionParts: ver,
		}
		v, b, err := intClient.ServerVersion.Get(context.Background())

		assert.NoError(t, err)
		assert.Nil(t, b)

		assert.Equal(t, &ServerVersion{
			Version:      "18.7.0",
			BuildNumber:  "7121",
			GitSha:       "75d1247f58ab8bcde3c5b43392a87347979f82c5",
			FullVersion:  "18.7.0 (7121-75d1247f58ab8bcde3c5b43392a87347979f82c5)",
			CommitURL:    "https://github.com/gocd/gocd/commits/75d1247f58ab8bcde3c5b43392a87347979f82c5",
			VersionParts: ver,
		}, v)
	}
}
