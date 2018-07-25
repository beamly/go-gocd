package gocd

import (
	"context"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServerVersion(t *testing.T) {
	t.Run("ServerVersionCaching", testServerVersionCaching)
	t.Run("Resource", testServerVersionResource)
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
