package gocd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipelineService(t *testing.T) {
	setup()
	defer teardown()

	t.Run("Get", testPipelineServiceGet)
	t.Run("Create/Delete", testPipelineServiceCreateDelete)
	t.Run("GetHistory", testPipelineServiceGetHistory)
	t.Run("GetStatus", testPipelineServiceGetStatus)
	t.Run("Un/Pause", testPipelineServiceUnPause)
	//t.Run("ReleaseLock", testPipelineServiceReleaseLock)
	t.Run("PaginationStub", testPipelineServicePaginationStub)
	t.Run("StageContainer", testPipelineStageContainer)
	t.Run("ConfirmHeader", testChoosePipelineConfirmHeader)
}

func TestPipelineServiceSchedule(t *testing.T) {
	for _, tt := range []struct {
		name        string
		versionFile string
	}{
		{
			name:        "unversionned",
			versionFile: "test/resources/version.0.json",
		},
		{
			name:        "18.2.0",
			versionFile: "test/resources/version.2.json",
		},
	} {

		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()

			mux.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
				j, _ := ioutil.ReadFile(tt.versionFile)
				fmt.Fprint(w, string(j))
			})

			mux.HandleFunc("/api/pipelines/test-pipeline/schedule", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, r.Method, "POST", "Unexpected HTTP method")
				fmt.Fprint(w, string([]byte(`{"message" : "Request to schedule pipeline test-pipeline accepted"}`)))
			})

			result, _, err := client.Pipelines.Schedule(context.Background(), "test-pipeline", nil)
			if err != nil {
				t.Error(err)
			}

			assert.True(t, result)
		})
	}

}

func TestPipelineServiceScheduleWithBody(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		j, _ := ioutil.ReadFile("test/resources/version.2.json")
		fmt.Fprint(w, string(j))
	})

	mux.HandleFunc("/api/pipelines/test-pipeline/schedule", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST", "Unexpected HTTP method")
		expReqBody, _ := ioutil.ReadFile("test/request/schedule-pipeline.0.json")
		reqBody, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, string(expReqBody), string(reqBody))
		fmt.Fprint(w, string([]byte(`{"message" : "Request to schedule pipeline test-pipeline accepted"}`)))
	})

	body := &ScheduleRequestBody{
		Materials: []*ScheduleMaterial{
			{
				Revision:    "123",
				Fingerprint: "45",
			},
			{
				Revision:    "67",
				Fingerprint: "89",
			},
		},
		EnvironmentVariables: []*EnvironmentVariable{
			{
				Name:   "USERNAME",
				Value:  "gocd",
				Secure: false,
			},
			{
				Name:   "SSH_PASSPHRASE",
				Value:  "some passphrase",
				Secure: true,
			},
			{
				Name:           "PASSWORD",
				EncryptedValue: "YEepp1G0C05SpP0fcp4Jh",
				Secure:         true,
			},
		},
	}

	result, _, err := client.Pipelines.Schedule(context.Background(), "test-pipeline", body)
	if err != nil {
		t.Error(err)
	}

	assert.True(t, result)
}

func TestPipelineServiceScheduleWithBodyForUnversionnedAPI(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Unexpected HTTP method")
		j, _ := ioutil.ReadFile("test/resources/version.0.json")
		fmt.Fprint(w, string(j))
	})

	mux.HandleFunc("/api/pipelines/test-pipeline/schedule", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Unexpected HTTP method")
		err := r.ParseForm()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, []string{"test"}, r.Form["materials[svn]"])
		assert.Equal(t, []string{"123"}, r.Form["materials[pkg-name]"])
		assert.Equal(t, []string{"gocd"}, r.Form["variables[USERNAME]"])
		assert.Equal(t, []string{"some passphrase"}, r.Form["secure_variables[SSH_PASSPHRASE]"])
		fmt.Fprint(w, string([]byte(`{"message" : "Request to schedule pipeline test-pipeline accepted"}`)))
	})

	body := &ScheduleRequestBody{
		Materials: []*ScheduleMaterial{
			{
				Name:     "svn",
				Revision: "test",
			},
			{
				Name:     "pkg-name",
				Revision: "123",
			},
		},
		EnvironmentVariables: []*EnvironmentVariable{
			{
				Name:   "USERNAME",
				Value:  "gocd",
				Secure: false,
			},
			{
				Name:   "SSH_PASSPHRASE",
				Value:  "some passphrase",
				Secure: true,
			},
		},
	}

	result, _, err := client.Pipelines.Schedule(context.Background(), "test-pipeline", body)
	if err != nil {
		t.Error(err)
	}

	assert.True(t, result)
}

func testPipelineStageContainer(t *testing.T) {

	p := &Pipeline{
		Name:   "mock-name",
		Stages: []*Stage{{Name: "1"}, {Name: "2"}},
	}

	i := StageContainer(p)

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

func testPipelineServicePaginationStub(t *testing.T) {
	pgs := PipelinesService{}

	assert.Equal(t, "a/b/c/4",
		pgs.buildPaginatedStub("a/%s/c", "b", 4))

	assert.Equal(t, "a/b/c",
		pgs.buildPaginatedStub("a/%s/c", "b", 0))

}

func testPipelineServiceGetStatus(t *testing.T) {
	mux.HandleFunc("/api/pipelines/test-pipeline/status", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		b, _ := ioutil.ReadFile("test/resources/pipeline.2.json")
		fmt.Fprint(w, string(b))
	})

	ps, _, err := client.Pipelines.GetStatus(context.Background(), "test-pipeline", 0)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, ps)
	assert.False(t, ps.Locked)
	assert.True(t, ps.Paused)
	assert.False(t, ps.Schedulable)
}

func testPipelineServiceCreateDelete(t *testing.T) {
	if !runIntegrationTest(t) {
		t.Skip("Skipping acceptance tests as GOCD_ACC not set to 1")
	}

	p := Pipeline{
		LabelTemplate:         "${COUNT}",
		EnablePipelineLocking: true,
		Name:                  "testPipelineServiceCreateDelete",
		Materials: []Material{
			{
				Type: "git",
				Attributes: MaterialAttributesGit{
					URL:          "git@github.com:sample_repo/example.git",
					Destination:  "dest",
					InvertFilter: false,
					AutoUpdate:   true,
					Branch:       "master",
					ShallowClone: true,
				},
			},
		},
		Stages: []*Stage{
			{
				Name:           "defaultStage",
				FetchMaterials: true,
				Approval: &Approval{
					Type: "success",
					Authorization: &Authorization{
						Roles: []string{},
						Users: []string{},
					},
				},
				Jobs: []*Job{
					{
						Name: "defaultJob",
						Tasks: []*Task{
							{
								Type: "exec",
								Attributes: TaskAttributes{
									RunIf:   []string{"passed"},
									Command: "ls",
								},
							},
						},
					},
				},
			},
		},
		Version: "mock-version",
	}

	ctx := context.Background()
	pr, _, err := intClient.PipelineConfigs.Create(ctx, "test-group", &p)
	assert.NoError(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, "testPipelineServiceCreateDelete", pr.Name)

	assert.Len(t, pr.Stages, 1)

	stage := pr.Stages[0]
	assert.Equal(t, "defaultStage", stage.Name)

	msg, _, err := intClient.PipelineConfigs.Delete(ctx, p.Name)
	assert.NoError(t, err)
	assert.Contains(t, msg, "'testPipelineServiceCreateDelete' was deleted successfully")
}

func testPipelineServiceGet(t *testing.T) {
	mux.HandleFunc("/api/pipelines/test-pipeline/instance/1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		j, _ := ioutil.ReadFile("test/resources/pipeline.0.json")
		fmt.Fprint(w, string(j))
	})

	p, _, err := client.Pipelines.GetInstance(context.Background(), "test-pipeline", 1)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, p)
	assert.Equal(t, p.Name, "test-pipeline")

	assert.Len(t, p.Stages, 1)

	s := p.Stages[0]
	assert.Equal(t, "stage1", s.Name)
}

func testPipelineServiceGetHistory(t *testing.T) {
	mux.HandleFunc("/api/pipelines/test-pipeline/history", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		j, _ := ioutil.ReadFile("test/resources/pipeline.1.json")
		fmt.Fprint(w, string(j))
	})
	ph, _, err := client.Pipelines.GetHistory(context.Background(), "test-pipeline", 0)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, ph)
	assert.Len(t, ph.Pipelines, 2)

	h1 := ph.Pipelines[0]
	assert.True(t, h1.CanRun)
	assert.Equal(t, h1.Name, "pipeline1")
	assert.Equal(t, h1.NaturalOrder, float32(11))
	assert.Equal(t, h1.Comment, "")
	assert.Equal(t, h1.Label, "11")
	assert.Equal(t, h1.Counter, 11)
	assert.Equal(t, h1.PreparingToSchedule, false)
	assert.Len(t, h1.Stages, 1)

	h1s := h1.Stages[0]
	assert.Equal(t, h1s.Name, "stage1")

	h2 := ph.Pipelines[1]
	assert.True(t, h2.CanRun)
	assert.Equal(t, h2.Name, "pipeline1")
	assert.Equal(t, h2.NaturalOrder, float32(10))
	assert.Equal(t, h2.Comment, "")
	assert.Equal(t, h2.Label, "10")
	assert.Equal(t, h2.Counter, 10)
	assert.Equal(t, h2.PreparingToSchedule, false)
	assert.Len(t, h2.Stages, 1)

	h2s := h2.Stages[0]
	assert.Equal(t, h2s.Name, "stage1")

}

func testChoosePipelineConfirmHeader(t *testing.T) {
	for _, tt := range []struct {
		name             string
		apiVersion       string
		wantHeaders      map[string]string
		wantResponseType string
		wantResponseBody interface{}
	}{
		{
			name:             "confirm",
			apiVersion:       "",
			wantHeaders:      map[string]string{"Confirm": "true"},
			wantResponseType: "",
			wantResponseBody: nil,
		},
		{
			name:             "x-confirm-v1",
			apiVersion:       "application/vnd.go.cd.v1+json",
			wantHeaders:      map[string]string{"X-GoCD-Confirm": "true"},
			wantResponseType: "json",
			wantResponseBody: &map[string]interface{}{},
		},
		{
			name:             "x-confirm-v2",
			apiVersion:       "application/vnd.go.cd.v2+json",
			wantHeaders:      map[string]string{"X-GoCD-Confirm": "true"},
			wantResponseType: "json",
			wantResponseBody: &map[string]interface{}{},
		},
		{
			name:             "x-confirm-v3",
			apiVersion:       "application/vnd.go.cd.v3+json",
			wantHeaders:      map[string]string{"X-GoCD-Confirm": "true"},
			wantResponseType: "json",
			wantResponseBody: &map[string]interface{}{},
		},
		{
			name:             "x-confirm-v4",
			apiVersion:       "application/vnd.go.cd.v4+json",
			wantHeaders:      map[string]string{"X-GoCD-Confirm": "true"},
			wantResponseType: "json",
			wantResponseBody: &map[string]interface{}{},
		},
		{
			name:             "x-confirm-v5",
			apiVersion:       "application/vnd.go.cd.v5+json",
			wantHeaders:      map[string]string{"X-GoCD-Confirm": "true"},
			wantResponseType: "json",
			wantResponseBody: &map[string]interface{}{},
		},
		{
			name:             "x-confirm-v6",
			apiVersion:       "application/vnd.go.cd.v6+json",
			wantHeaders:      map[string]string{"X-GoCD-Confirm": "true"},
			wantResponseType: "json",
			wantResponseBody: &map[string]interface{}{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			request := &APIClientRequest{}
			choosePipelineConfirmHeader(request, tt.apiVersion)
			assert.Equal(t, tt.wantHeaders, request.Headers)
			assert.Equal(t, tt.wantResponseBody, request.ResponseBody)
			assert.Equal(t, tt.wantResponseType, request.ResponseType)
		})
	}
}
