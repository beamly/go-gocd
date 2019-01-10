package gocd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestPipelineTemplate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		j, _ := ioutil.ReadFile("test/resources/version.1.json")
		fmt.Fprint(w, string(j))
	})

	t.Run("List", testListPipelineTemplates)
	t.Run("Get", testGetPipelineTemplate)
	t.Run("Delete", testDeletePipelineTemplate)
	t.Run("RemoveLinks", tesPipelineTemplateRemoveLinks)
	t.Run("Pipelines", testPipelineTemplatePipelines)
	t.Run("Update", testPipelineTemplateUpdate)
	t.Run("StageContainerI", testPipelineTemplateStageContainer)
}

func testPipelineTemplateStageContainer(t *testing.T) {
	var i StageContainer

	i = &PipelineTemplate{
		Name:   "mock-name",
		Stages: []*Stage{{Name: "1"}, {Name: "2"}},
	}

	assert.Equal(t, "mock-name", i.GetName())
	assert.Len(t, i.GetStages(), 2)

	i.AddStage(&Stage{Name: "3"})
	assert.Len(t, i.GetStages(), 3)

	s1 := i.GetStage("1")
	assert.Equal(t, s1.Name, "1")

	sn := i.GetStage("hello")
	assert.Nil(t, sn)

	i.SetStages([]*Stage{})
	assert.Len(t, i.GetStages(), 0)
}

// TestPipelineTemplateCreate is a seperate test to avoid overlapping mock HandleFunc's.
func TestPipelineTemplateCreate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/admin/templates", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST", "Unexpected HTTP method")
		apiVersion, _ := client.getAPIVersion(context.Background(), "admin/templates")
		assert.Equal(t, apiVersion, r.Header.Get("Accept"))

		j, _ := ioutil.ReadFile("test/resources/pipelinetemplate.2.json")
		w.Header().Set("Etag", "mock-etag")

		fmt.Fprint(w, string(j))
	})
	pt, _, err := client.PipelineTemplates.Create(context.Background(),
		"test-config2",
		[]*Stage{{}},
	)
	assert.NoError(t, err)
	assert.NotNil(t, pt)
	assert.Equal(t, "mock-etag", pt.Version)
}

func testPipelineTemplateUpdate(t *testing.T) {
	mux.HandleFunc("/api/admin/templates/test-config", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT", "Unexpected HTTP method")
		j, _ := ioutil.ReadFile("test/resources/pipelinetemplate.1.json")
		fmt.Fprint(w, string(j))
	})

	pt, _, err := client.PipelineTemplates.Update(context.Background(),
		"test-config",
		&PipelineTemplate{
			Stages: []*Stage{
				{},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, pt)

	assert.Equal(t, "template", pt.Name)
	assert.Len(t, pt.Stages, 1)

	s := pt.Stages[0]
	assert.Equal(t, "defaultStage", s.Name)
	assert.True(t, s.FetchMaterials)
	assert.False(t, s.CleanWorkingDirectory)
	assert.False(t, s.NeverCleanupArtifacts)

	assert.Len(t, s.EnvironmentVariables, 0)
	assert.Len(t, s.Resources, 0)
}

func testPipelineTemplatePipelines(t *testing.T) {
	p := []*Pipeline{}
	pt := PipelineTemplate{Embedded: &embeddedPipelineTemplate{Pipelines: p}}

	assert.Exactly(t, p, pt.Pipelines())
}

func tesPipelineTemplateRemoveLinks(t *testing.T) {
	pt := PipelineTemplate{Links: &HALLinks{}}
	assert.NotNil(t, pt.Links)
	pt.RemoveLinks()
	assert.Nil(t, pt.Links)
}

func testDeletePipelineTemplate(t *testing.T) {

	b, err := json.Marshal(map[string]string{
		"message": "The template 'template2' was deleted successfully.",
	})
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/api/admin/templates/template2", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE", "Unexpected HTTP method")
		fmt.Fprint(w, string(b))
	})

	message, _, err := client.PipelineTemplates.Delete(context.Background(), "template2")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "The template 'template2' was deleted successfully.", message)

}

func testListPipelineTemplates(t *testing.T) {

	mux.HandleFunc("/api/admin/templates", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		testAuth(t, r, mockAuthorization)
		j, _ := ioutil.ReadFile("test/resources/pipelinetemplates.0.json")
		fmt.Fprint(w, string(j))
	})

	templates, _, err := client.PipelineTemplates.List(context.Background())

	assert.Nil(t, err)
	assert.Len(t, templates, 1)
	for _, attribute := range []EqualityTest{
		{templates[0].Name, "template0"},
		{templates[0].Embedded.Pipelines[0].Name, "up42"},
	} {
		assert.Equal(t, attribute.got, attribute.wanted)
	}
}

func testGetPipelineTemplate(t *testing.T) {

	mux.HandleFunc("/api/admin/templates/template1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		testAuth(t, r, mockAuthorization)
		j, _ := ioutil.ReadFile("test/resources/pipelinetemplate.0.json")
		w.Header().Set("Etag", "mock-etag")
		fmt.Fprint(w, string(j))
	})

	template, _, err := client.PipelineTemplates.Get(
		context.Background(),
		"template1",
	)

	assert.NoError(t, err)
	assert.Len(t, template.Stages, 1)

	assert.Equal(t, "mock-etag", template.Version)

	for _, attribute := range []EqualityTest{
		{template.Name, "template1"},
		{template.Stages[0].Name, "up42_stage"},
		{template.Stages[0].Approval.Type, "success"},
	} {
		assert.Equal(t, attribute.got, attribute.wanted)
	}
}
