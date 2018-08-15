package gocd

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestPluginApi(t *testing.T) {
	setup()
	defer teardown()

	t.Run("List", testPluginAPIList)
	t.Run("Get", testPluginAPIGet)
}

func testPluginAPIList(t *testing.T) {
	cachedServerVersion = &ServerVersion{Version: "17.1.0"}
	err := cachedServerVersion.parseVersion()
	assert.NoError(t, err)

	mux.HandleFunc("/api/admin/plugin_info", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Unexpected HTTP method")
		j, _ := ioutil.ReadFile("test/resources/plugin.0.json")
		fmt.Fprint(w, string(j))
	})

	plugins, _, err := client.Plugins.List(context.Background())
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, plugins)

	assert.NotNil(t, plugins.Links.Get("Doc"))
	assert.Equal(t, "https://api.gocd.org/#plugin-info", plugins.Links.Get("Doc").URL.String())
	assert.NotNil(t, plugins.Links.Get("Self"))
	assert.Equal(t, "https://ci.example.com/go/api/admin/plugin_info", plugins.Links.Get("Self").URL.String())

	assert.NotNil(t, plugins.Embedded)
	assert.NotNil(t, plugins.Embedded.PluginInfo)
	assert.Len(t, plugins.Embedded.PluginInfo, 1)

	pi := plugins.Embedded.PluginInfo[0]
	assert.NotNil(t, pi.Links)
	assert.Equal(t, "https://ci.example.com/go/api/admin/plugin_info/plugin_id", pi.Links.Get("Self").URL.String())
	assert.Equal(t, "https://api.gocd.org/#plugin-info", pi.Links.Get("Doc").URL.String())
	assert.Equal(t, "https://ci.example.com/go/api/admin/plugin_info/:id", pi.Links.Get("Find").URL.String())

	assert.Equal(t, "test-plugin", pi.ID)
	assert.Equal(t, "SCM Plugin", pi.Name)
	assert.Equal(t, "SCM Plugin For HG", pi.DisplayName)
	assert.Equal(t, "1.2.3", pi.Version)
	assert.Equal(t, "scm", pi.Type)
}

func testPluginAPIGet(t *testing.T) {
	cachedServerVersion = &ServerVersion{Version: "17.1.0"}
	err := cachedServerVersion.parseVersion()
	assert.NoError(t, err)

	mux.HandleFunc("/api/admin/plugin_info/test-plugin", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Unexpected HTTP method")
		j, _ := ioutil.ReadFile("test/resources/plugin.1.json")
		fmt.Fprint(w, string(j))
	})

	plugin, _, err := client.Plugins.Get(context.Background(), "test-plugin")
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, plugin)

	assert.Equal(t, "test-plugin", plugin.ID)
	assert.Equal(t, "My SCM Plugin", plugin.Name)
	assert.Equal(t, "SCM Plugin For HG", plugin.DisplayName)
	assert.Equal(t, "1.2.3", plugin.Version)
	assert.Equal(t, "scm", plugin.Type)
	assert.NotNil(t, plugin.PluggableInstanceSettings)

	assert.Len(t, plugin.PluggableInstanceSettings.Configurations, 3)

	cfg0 := plugin.PluggableInstanceSettings.Configurations[0]
	assert.Equal(t, "url", cfg0.Key)
	assert.False(t, cfg0.Metadata.Secure)
	assert.True(t, cfg0.Metadata.Required)
	assert.True(t, cfg0.Metadata.PartOfIdentity)

	cfg1 := plugin.PluggableInstanceSettings.Configurations[1]
	assert.Equal(t, "username", cfg1.Key)
	assert.False(t, cfg1.Metadata.Secure)
	assert.False(t, cfg1.Metadata.Required)
	assert.False(t, cfg1.Metadata.PartOfIdentity)

	cfg2 := plugin.PluggableInstanceSettings.Configurations[2]
	assert.Equal(t, "password", cfg2.Key)
	assert.True(t, cfg2.Metadata.Secure)
	assert.False(t, cfg2.Metadata.Required)
	assert.False(t, cfg2.Metadata.PartOfIdentity)

	assert.Equal(t,
		"<div>Plugin view template </div>",
		plugin.PluggableInstanceSettings.View.Template,
	)

}
