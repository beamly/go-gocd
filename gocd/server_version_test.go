package gocd

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServerVersion(t *testing.T) {
	t.Run("ServerVersion", testServerVersion)
	t.Run("BadServerVersion", testBadServerVersion)
	t.Run("ServerVersionCaching", testServerVersionCaching)
	t.Run("Resource", testServerVersionResource)
}

func testServerVersion(t *testing.T) {

	if runIntegrationTest() {

		cachedServerVersion = nil
		v, _, err := client.ServerVersion.Get(context.Background())

		assert.NoError(t, err)

		assert.Equal(t, &ServerVersion{
			Version:     "16.6.0",
			BuildNumber: "3348",
			GitSha:      "a7a5717cbd60c30006314fb8dd529796c93adaf0",
			FullVersion: "16.6.0 (3348-a7a5717cbd60c30006314fb8dd529796c93adaf0)",
			CommitURL:   "https://github.com/gocd/gocd/commits/a7a5717cbd60c30006314fb8dd529796c93adaf0",
			VersionParts: &ServerVersionParts{
				Major: 16,
				Minor: 6,
				Patch: 0,
			},
		}, v)

		// Verify that the server version is cached
		assert.Equal(t, cachedServerVersion, v)
	}
}

func testBadServerVersion(t *testing.T) {
	for _, test := range []struct {
		name      string
		id        int
		errString string
	}{
		{name: "Major", id: 2, errString: "strconv.Atoi: parsing \"a\": invalid syntax"},
		{name: "Minor", id: 3, errString: "strconv.Atoi: parsing \"b\": invalid syntax"},
		{name: "Patch", id: 4, errString: "strconv.Atoi: parsing \"c\": invalid syntax"},
	} {
		cachedServerVersion = nil
		t.Run(test.name, func(t *testing.T) { testBadServerVersionMajor(t, test.id, test.errString) })
	}
}

func testBadServerVersionMajor(t *testing.T, i int, errString string) {
	if runIntegrationTest() {
		_, _, err := client.ServerVersion.Get(context.Background())
		assert.EqualError(t, err, errString)
	}
}

func testServerVersionCaching(t *testing.T) {

	cachedServerVersion = &ServerVersion{
		Version:     "18.7.0",
		BuildNumber: "7121",
		GitSha:      "75d1247f58ab8bcde3c5b43392a87347979f82c5",
		FullVersion: "18.7.0 (7121-75d1247f58ab8bcde3c5b43392a87347979f82c5)",
		CommitURL:   "https://github.com/gocd/gocd/commits/75d1247f58ab8bcde3c5b43392a87347979f82c5",
		VersionParts: &ServerVersionParts{
			Major: 18,
			Minor: 7,
			Patch: 0,
		},
	}
	v, b, err := client.ServerVersion.Get(context.Background())

	assert.NoError(t, err)
	assert.Nil(t, b)

	assert.Equal(t, &ServerVersion{
		Version:     "18.7.0",
		BuildNumber: "7121",
		GitSha:      "75d1247f58ab8bcde3c5b43392a87347979f82c5",
		FullVersion: "18.7.0 (7121-75d1247f58ab8bcde3c5b43392a87347979f82c5)",
		CommitURL:   "https://github.com/gocd/gocd/commits/75d1247f58ab8bcde3c5b43392a87347979f82c5",
		VersionParts: &ServerVersionParts{
			Major: 18,
			Minor: 7,
			Patch: 0,
		},
	}, v)
}
