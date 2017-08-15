package gocd

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestPipelineGroupsService(t *testing.T) {
	t.Run("List", testPipelineGroupsServiceList)
	t.Run("Filter", testPipelineGroupsServiceFilter)
}

func testPipelineGroupsServiceFilter(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/config/pipeline_groups", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		j, _ := ioutil.ReadFile("test/resources/pipelinegroups.1.json")
		fmt.Fprint(w, string(j))
	})
	pgs, _, err := client.PipelineGroups.List(context.Background(), "filter-group")
	assert.Nil(t, err)
	assert.Len(t, pgs, 1)

	pg := pgs[0]
	assert.Equal(t, "filter-group", pg.Name)

}

func testPipelineGroupsServiceList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/config/pipeline_groups", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		j, _ := ioutil.ReadFile("test/resources/pipelinegroups.0.json")
		fmt.Fprint(w, string(j))
	})
	pgs, _, err := client.PipelineGroups.List(context.Background(), "")

	assert.Nil(t, err)
	assert.Len(t, pgs, 1)

	pg := pgs[0]
	assert.Equal(t, pg.Name, "first")

	assert.Len(t, pg.Pipelines, 1)

	p := pg.Pipelines[0]
	assert.Equal(t, p.Name, "up42")
	assert.Equal(t, p.Label, "${COUNT}")
	assert.Len(t, p.Stages, 1)
	assert.Len(t, p.Materials, 1)

	s := p.Stages[0]
	assert.Equal(t, s.Name, "up42_stage")

	m := p.Materials[0]
	assert.Equal(t, m.Type, "Git")
	assert.Equal(t, m.Fingerprint, "2d05446cd52a998fe3afd840fc2c46b7c7e421051f0209c7f619c95bedc28b88")
	assert.Equal(t, m.Description, "URL: https://github.com/gocd/gocd, Branch: master")
}
