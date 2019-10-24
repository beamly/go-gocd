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
	t.Run("Integration", testEnvironmentIntegration)

	setup()
	defer teardown()

	mux.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		j, _ := ioutil.ReadFile("test/resources/version.3.json")
		fmt.Fprint(w, string(j))
	})

	t.Run("List", testEnvironmentList)
	t.Run("Delete", testEnvironmentDelete)
	t.Run("Get", testEnvironmentGet)
	t.Run("Patch", testEnvironmentPatch)
}

func testEnvironmentIntegration(t *testing.T) {
	if !runIntegrationTest(t) {
		t.Skip("Skipping acceptance tests as GOCD_ACC not set to 1")
	}

	ctx := context.Background()

	env, _, err := intClient.Environments.Create(ctx, "test")
	if err != nil {
		t.Error(err)
	}

	p := &Pipeline{
		Name: "environment-pipeline",
		Materials: []Material{{
			Type: "git",
			Attributes: MaterialAttributesGit{
				URL:         "git@github.com:sample_repo/example.git",
				Destination: "dest",
				Branch:      "master",
			},
		}},
		Stages: buildUpstreamPipelineStages(),
	}

	_, _, err = intClient.PipelineConfigs.Create(ctx, "test-group", p)
	if err != nil {
		t.Error(err)
	}

	patch := EnvironmentPatchRequest{
		Pipelines: &PatchStringAction{
			Add:    []string{"environment-pipeline"},
			Remove: []string{},
		},
		EnvironmentVariables: &EnvironmentVariablesAction{
			Add: []*EnvironmentVariable{
				{
					Name:  "GO_SERVER_URL",
					Value: "https://ci.example.com/go",
				},
			},
			Remove: []string{},
		},
	}
	_, _, err = intClient.Environments.Patch(context.Background(), "test", &patch)
	if err != nil {
		t.Error(err)
	}

	envs, _, err := intClient.Environments.List(ctx)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, envs)

	// Make sure version-specific expected values are set
	apiVersion, err := intClient.getAPIVersion(ctx, "admin/environments")
	assert.NoError(t, err)

	var envDoc, envFind, pipelineDoc, pipelineFind, pipelineSelf string

	switch apiVersion {
	case apiV3:
		envDoc = "https://api.go.cd/current/#environment-config"
		envFind = "http://127.0.0.1:8153/go/api/admin/environments/:name"
		pipelineDoc = "https://api.go.cd/current/#pipelines"
		pipelineFind = "/api/admin/pipelines/:pipeline_name"
		pipelineSelf = "http://127.0.0.1:8153/go/api/pipelines/environment-pipeline/history"
	case apiV2:
		envDoc = "https://api.gocd.org/#environment-config"
		envFind = "http://127.0.0.1:8153/go/api/admin/environments/:environment_name"
		pipelineDoc = "https://api.gocd.org/#pipeline-config"
		pipelineFind = "http://127.0.0.1:8153/go/api/admin/pipelines/:pipeline_name"
		pipelineSelf = "http://127.0.0.1:8153/go/api/admin/pipelines/environment-pipeline"
	}

	assert.NotNil(t, envs.Links.Get("Self"))
	assert.Equal(t, "http://127.0.0.1:8153/go/api/admin/environments", envs.Links.Get("Self").URL.String())
	assert.NotNil(t, envs.Links.Get("Doc"))
	assert.Equal(t, envDoc, envs.Links.Get("Doc").URL.String())

	assert.NotNil(t, envs.Embedded)
	assert.NotNil(t, envs.Embedded.Environments)
	assert.Len(t, envs.Embedded.Environments, 1)

	env = envs.Embedded.Environments[0]
	assert.NotNil(t, env.Links)
	assert.Equal(t, "http://127.0.0.1:8153/go/api/admin/environments/test", env.Links.Get("Self").URL.String())
	assert.Equal(t, envDoc, env.Links.Get("Doc").URL.String())
	assert.Equal(t, envFind, env.Links.Get("Find").URL.String())

	assert.Equal(t, "test", env.Name)

	assert.NotNil(t, env.Pipelines)
	assert.Len(t, env.Pipelines, 1)

	p = env.Pipelines[0]
	assert.NotNil(t, p.Links)
	assert.Equal(t, pipelineSelf, p.Links.Get("Self").URL.String())
	assert.Equal(t, pipelineDoc, p.Links.Get("Doc").URL.String())
	assert.Equal(t, pipelineFind, p.Links.Get("Find").URL.String())
	assert.Equal(t, "environment-pipeline", p.Name)

	assert.NotNil(t, env.EnvironmentVariables)
	assert.Len(t, env.EnvironmentVariables, 1)

	ev1 := env.EnvironmentVariables[0]
	assert.Equal(t, "GO_SERVER_URL", ev1.Name)
	assert.False(t, ev1.Secure)
	assert.Equal(t, "https://ci.example.com/go", ev1.Value)
}

func testEnvironmentList(t *testing.T) {
	apiVersion, err := client.getAPIVersion(context.Background(), "admin/environments")
	assert.NoError(t, err)

	mux.HandleFunc("/api/admin/environments", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		assert.Contains(t, r.Header["Accept"], apiVersion)

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
	apiVersion, err := client.getAPIVersion(context.Background(), "admin/environments/:environment_name")
	assert.NoError(t, err)

	mux.HandleFunc("/api/admin/environments/my_environment_1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE", "Unexpected HTTP method")
		assert.Contains(t, r.Header["Accept"], apiVersion)

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
	apiVersion, err := client.getAPIVersion(context.Background(), "admin/environments/:environment_name")
	assert.NoError(t, err)

	mux.HandleFunc("/api/admin/environments/my_environment", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		assert.Contains(t, r.Header["Accept"], apiVersion)

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
	apiVersion, err := client.getAPIVersion(context.Background(), "admin/environments/:environment_name")
	assert.NoError(t, err)

	mux.HandleFunc("/api/admin/environments/my_environment_2", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PATCH", "Unexpected HTTP method")
		assert.Contains(t, r.Header["Accept"], apiVersion)

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
