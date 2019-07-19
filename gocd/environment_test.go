package gocd

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestEnvironment(t *testing.T) {
	setup()
	defer teardown()

	t.Run("List", testEnvironmentList)
	t.Run("Delete", testEnvironmentDelete)
	t.Run("Get", testEnvironmentGet)
	t.Run("Patch", testEnvironmentPatch)
}

func testEnvironmentList(t *testing.T) {
	mux.HandleFunc("/api/admin/environments", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		assert.Contains(t, r.Header["Accept"], "application/vnd.go.cd.v2+json")

		j, _ := ioutil.ReadFile("test/resources/environment.0.json")
		fmt.Fprint(w, string(j))
	})

	envs, _, err := client.Environments.List(context.Background())
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, envs)

	assert.NotNil(t, envs.Links.Get("Self"))
	assert.Equal(t, "https://ci.example.com/go/api/admin/environments", envs.Links.Get("Self").URL.String())
	assert.NotNil(t, envs.Links.Get("Doc"))
	assert.Equal(t, "https://api.gocd.org/#environment-config", envs.Links.Get("Doc").URL.String())

	assert.NotNil(t, envs.Embedded)
	assert.NotNil(t, envs.Embedded.Environments)
	assert.Len(t, envs.Embedded.Environments, 1)

	env := envs.Embedded.Environments[0]
	assert.NotNil(t, env.Links)
	assert.Equal(t, "https://ci.example.com/go/api/admin/environments/foobar", env.Links.Get("Self").URL.String())
	assert.Equal(t, "https://api.gocd.org/#environment-config", env.Links.Get("Doc").URL.String())
	assert.Equal(t, "https://ci.example.com/go/api/admin/environments/:environment_name", env.Links.Get("Find").URL.String())

	assert.Equal(t, "foobar", env.Name)

	assert.NotNil(t, env.Pipelines)
	assert.Len(t, env.Pipelines, 1)

	p := env.Pipelines[0]
	assert.NotNil(t, p.Links)
	assert.Equal(t, "https://ci.example.com/go/api/admin/pipelines/up42", p.Links.Get("Self").URL.String())
	assert.Equal(t, "https://api.gocd.org/#pipeline-config", p.Links.Get("Doc").URL.String())
	assert.Equal(t, "https://ci.example.com/go/api/admin/pipelines/:pipeline_name", p.Links.Get("Find").URL.String())
	assert.Equal(t, "up42", p.Name)

	assert.NotNil(t, env.Agents)
	assert.Len(t, env.Agents, 1)

	a := env.Agents[0]
	assert.NotNil(t, a.Links)
	assert.Equal(t, "https://ci.example.com/go/api/agents/adb9540a-b954-4571-9d9b-2f330739d4da", a.Links.Get("Self").URL.String())
	assert.Equal(t, "https://api.gocd.org/#agents", a.Links.Get("Doc").URL.String())
	assert.Equal(t, "https://ci.example.com/go/api/agents/:uuid", a.Links.Get("Find").URL.String())
	assert.Equal(t, "12345678-e2f6-4c78-123456789012", a.UUID)

	assert.NotNil(t, env.EnvironmentVariables)
	assert.Len(t, env.EnvironmentVariables, 2)

	ev1 := env.EnvironmentVariables[0]
	assert.Equal(t, "username", ev1.Name)
	assert.False(t, ev1.Secure)
	assert.Equal(t, "admin", ev1.Value)

	ev2 := env.EnvironmentVariables[1]
	assert.Equal(t, "password", ev2.Name)
	assert.True(t, ev2.Secure)
	assert.Equal(t, "LSd1TI0eLa+DjytHjj0qjA==", ev2.EncryptedValue)
}

func testEnvironmentDelete(t *testing.T) {
	mux.HandleFunc("/api/admin/environments/my_environment_1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE", "Unexpected HTTP method")
		assert.Contains(t, r.Header["Accept"], "application/vnd.go.cd.v2+json")

		fmt.Fprint(w, `{
  "message": "Environment 'my_environment_1' was deleted successfully."
}`)
	})

	message, resp, err := client.Environments.Delete(context.Background(), "my_environment_1")
	if err != nil {
		assert.Error(t, err)
	}
	assert.NotNil(t, resp)
	assert.Equal(t, "Environment 'my_environment_1' was deleted successfully.", message)
}

func testEnvironmentGet(t *testing.T) {
	mux.HandleFunc("/api/admin/environments/my_environment", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		assert.Contains(t, r.Header["Accept"], "application/vnd.go.cd.v2+json")

		j, _ := ioutil.ReadFile("test/resources/environment.1.json")
		fmt.Fprint(w, string(j))
	})

	env, _, err := client.Environments.Get(context.Background(), "my_environment")
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, env)

	assert.Equal(t, "https://ci.example.com/go/api/admin/environments/my_environment", env.Links.Get("self").URL.String())
	assert.Equal(t, "https://api.gocd.org/#environment-config", env.Links.Get("doc").URL.String())
	assert.Equal(t, "https://ci.example.com/go/api/admin/environments/:environment_name", env.Links.Get("find").URL.String())

	assert.Equal(t, "my_environment", env.Name)

	assert.Len(t, env.Pipelines, 1)
	p := env.Pipelines[0]
	assert.Equal(t, "https://ci.example.com/go/api/admin/pipelines/up42", p.Links.Get("self").URL.String())
	assert.Equal(t, "https://api.gocd.org/#pipeline-config", p.Links.Get("doc").URL.String())
	assert.Equal(t, "https://ci.example.com/go/api/admin/pipelines/:pipeline_name", p.Links.Get("find").URL.String())
	assert.Equal(t, "up42", p.Name)

	assert.Len(t, env.Agents, 1)
	a := env.Agents[0]
	assert.Equal(t, "https://ci.example.com/go/api/agents/adb9540a-b954-4571-9d9b-2f330739d4da", a.Links.Get("self").URL.String())
	assert.Equal(t, "https://api.gocd.org/#agents", a.Links.Get("doc").URL.String())
	assert.Equal(t, "https://ci.example.com/go/api/agents/:uuid", a.Links.Get("find").URL.String())
	assert.Equal(t, "12345678-e2f6-4c78-123456789012", a.UUID)

	assert.Len(t, env.EnvironmentVariables, 2)
	assert.Equal(t,
		&EnvironmentVariable{
			Secure: false,
			Name:   "username",
			Value:  "admin",
		},
		env.EnvironmentVariables[0],
	)

	assert.Equal(t,
		&EnvironmentVariable{
			Secure:         true,
			Name:           "password",
			EncryptedValue: "LSd1TI0eLa+DjytHjj0qjA==",
		},
		env.EnvironmentVariables[1],
	)

}

func testEnvironmentPatch(t *testing.T) {
	mux.HandleFunc("/api/admin/environments/my_environment_2", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PATCH", "Unexpected HTTP method")
		assert.Contains(t, r.Header["Accept"], "application/vnd.go.cd.v2+json")

		j, _ := ioutil.ReadFile("test/resources/environment.2.json")
		fmt.Fprint(w, string(j))

	})

	patch := EnvironmentPatchRequest{
		Pipelines: &PatchStringAction{
			Add:    []string{"up42"},
			Remove: []string{"sample"},
		},
		Agents: &PatchStringAction{
			Add:    []string{"12345678-e2f6-4c78-123456789012"},
			Remove: []string{"87654321-e2f6-4c78-123456789012"},
		},
		EnvironmentVariables: &EnvironmentVariablesAction{
			Add: []*EnvironmentVariable{
				{
					Name:  "GO_SERVER_URL",
					Value: "https://ci.example.com/go",
				},
			},
			Remove: []string{
				"URL",
			},
		},
	}
	env, _, err := client.Environments.Patch(context.Background(), "my_environment_2", &patch)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "https://ci.example.com/go/api/admin/environments/new_environment", env.Links.Get("self").URL.String())
	assert.Equal(t, "https://api.gocd.org/#environment-config", env.Links.Get("doc").URL.String())

	assert.Equal(t, "new_environment", env.Name)

	assert.Len(t, env.Pipelines, 1)
	p := env.Pipelines[0]
	assert.Equal(t, "https://ci.example.com/go/api/admin/pipelines/pipeline1", p.Links.Get("self").URL.String())
	assert.Equal(t, "https://api.gocd.org/#pipeline-config", p.Links.Get("doc").URL.String())
	assert.Equal(t, "https://ci.example.com/go/api/admin/pipelines/:pipeline_name", p.Links.Get("find").URL.String())
	assert.Equal(t, "up42", p.Name)

	assert.Len(t, env.Agents, 1)
	a := env.Agents[0]
	assert.Equal(t, "https://ci.example.com/go/api/agents/adb9540a-b954-4571-9d9b-2f330739d4da", a.Links.Get("self").URL.String())
	assert.Equal(t, "https://api.gocd.org/#agents", a.Links.Get("doc").URL.String())
	assert.Equal(t, "https://ci.example.com/go/api/agents/:uuid", a.Links.Get("find").URL.String())
	assert.Equal(t, "12345678-e2f6-4c78-123456789012", a.UUID)

	assert.Equal(t, []*EnvironmentVariable{
		{
			Secure: false,
			Name:   "GO_SERVER_URL",
			Value:  "https://ci.example.com/go",
		},
	}, env.EnvironmentVariables)

}
