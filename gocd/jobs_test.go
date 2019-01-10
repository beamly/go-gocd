package gocd

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestTaskValidate(t *testing.T) {
	setup()
	defer teardown()

	t.Run("ListScheduled", taskTaskValidateListScheduled)
	t.Run("Fail", taskValidateFail)
	t.Run("SuccessExec", taskValidateSuccessExec)
	t.Run("SuccessAnt", taskValidateSuccessAnt)
	t.Run("JSONString", testJobJSONString)
	t.Run("JSONStringFail", testJobJSONStringFail)
	t.Run("EmptyEnvironmentVariableValue", testEmptyEnvironmentVariableValue)
}

func testJobJSONStringFail(t *testing.T) {
	jb := Job{}
	_, err := jb.JSONString()
	assert.EqualError(t, err, "`gocd.Jobs.Name` is empty")
}

func testJobJSONString(t *testing.T) {
	jb := Job{Name: "test-job"}

	j, err := jb.JSONString()
	if err != nil {
		assert.Nil(t, err)
	}
	assert.Equal(
		t, `{
  "name": "test-job"
}`, j)
}

func testEmptyEnvironmentVariableValue(t *testing.T) {
	jb := Job{
		Name: "test-job",
		EnvironmentVariables: []*EnvironmentVariable{
			{
				Name:  "test",
				Value: "",
			},
		},
	}

	j, err := jb.JSONString()
	if err != nil {
		assert.Nil(t, err)
	}
	assert.Equal(
		t, `{
  "name": "test-job",
  "environment_variables": [
    {
      "name": "test",
      "value": "",
      "secure": false
    }
  ]
}`, j)
}

func taskTaskValidateListScheduled(t *testing.T) {
	mux.HandleFunc("/api/jobs/scheduled.xml", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET", "Unexpected HTTP method")
		j, _ := ioutil.ReadFile("test/resources/jobs.0.xml")
		fmt.Fprint(w, string(j))
	})
	sj, _, err := client.Jobs.ListScheduled(context.Background())
	if err != nil {
		assert.Nil(t, err)
	}
	assert.NotNil(t, sj)

	assert.Len(t, sj, 2)

	j1 := sj[0]
	assert.Equal(t, j1.Name, "job1")
	assert.Equal(t, j1.ID, "6")
	assert.NotNil(t, j1.Link)
	assert.Equal(t, j1.Link.HRef, "https://ci.example.com/go/tab/build/detail/mypipeline/5/defaultStage/1/job1")
	assert.Equal(t, j1.Link.Rel, "self")
	assert.Equal(t, j1.BuildLocator, "mypipeline/5/defaultStage/1/job1")

	j2 := sj[1]
	assert.Equal(t, j2.Name, "job2")
	assert.Equal(t, j2.ID, "7")
	assert.NotNil(t, j2.Link)
	assert.Equal(t, j2.Link.HRef, "https://ci.example.com/go/tab/build/detail/mypipeline/5/defaultStage/1/job2")
	assert.Equal(t, j2.Link.Rel, "self")
	assert.Equal(t, j2.BuildLocator, "mypipeline/5/defaultStage/1/job2")

}

func taskValidateSuccessAnt(t *testing.T) {
	antTask := Task{
		Type: "ant",
	}
	assert.NotNil(t, antTask.Validate())

	antTask.Attributes.RunIf = []string{"one", "two"}
	assert.NotNil(t, antTask.Validate())

	antTask.Attributes.BuildFile = "build-file"
	assert.NotNil(t, antTask.Validate())

	antTask.Attributes.Target = "target"
	assert.NotNil(t, antTask.Validate())

	antTask.Attributes.WorkingDirectory = "working-directory"
	assert.Nil(t, antTask.Validate())
}

func taskValidateSuccessExec(t *testing.T) {
	execTask := Task{
		Type: "exec",
	}
	assert.NotNil(t, execTask.Validate())

	execTask.Attributes.RunIf = []string{"one", "two"}
	assert.NotNil(t, execTask.Validate())

	execTask.Attributes.Command = "command-one"
	assert.NotNil(t, execTask.Validate())

	execTask.Attributes.Arguments = []string{"one", "two"}
	assert.NotNil(t, execTask.Validate())

	execTask.Attributes.WorkingDirectory = "one-two-three"
	assert.Nil(t, execTask.Validate())

}

func taskValidateFail(t *testing.T) {
	task := Task{}
	assert.EqualError(t,
		task.Validate(), "Missing `gocd.TaskAttribute` type")

	task.Type = "invalid-task-type"
	assert.EqualError(t,
		task.Validate(), "Unexpected `gocd.Task.Attribute` types")

	task.Type = "exec"
	assert.NotNil(t, task.Validate())

	task.Type = "ant"
	assert.NotNil(t, task.Validate())
}

func TestJobValidate(t *testing.T) {
	t.Run("ValidateJob", jobValidateSuccess)
	t.Run("Exec", jobValidateExecSuccess)
	t.Run("Ant", jobValidateAntSuccess)
	//t.Run("Nant", job_ValidateNantSuccess)
	//t.Run("Rake", job_ValidateRakeSuccess)
	//t.Run("Fetch", job_ValidateFetchSuccess)
	//t.Run("PluggableTask", job_ValidatePluggableTaskSuccess)
}

func jobValidateSuccess(t *testing.T) {
	j := Job{}
	err := j.Validate()
	assert.NotNil(t, err)

	j.Name = "job-name"
	err = j.Validate()
	assert.Nil(t, err)
}

func jobValidateExecSuccess(t *testing.T) {
	err := (&TaskAttributes{
		RunIf:            []string{"runif-exec"},
		Command:          "my-test-command",
		Arguments:        []string{"arg1", "arg2"},
		WorkingDirectory: "test-working-diretory",
	}).ValidateExec()
	assert.Nil(t, err)
}

func jobValidateAntSuccess(t *testing.T) {
	err := (&TaskAttributes{
		RunIf:            []string{"runif-ant"},
		BuildFile:        "test-build-file",
		Target:           "test-target",
		WorkingDirectory: "test-working-directory",
	}).ValidateAnt()
	assert.Nil(t, err)
}

//func job_ValidateNantSuccess(t *testing.T) {
//	err := (&TaskAttributes{}).ValidateNant()
//	assert.Nil(t, err)
//
//}
//
//func job_ValidateRakeSuccess(t *testing.T) {
//	err := (&TaskAttributes{}).ValidateRake()
//	assert.Nil(t, err)
//
//}
//
//func job_ValidateFetchSuccess(t *testing.T) {
//	err := (&TaskAttributes{}).ValidateFetch()
//	assert.Nil(t, err)
//
//}
//
//func job_ValidatePluggableTaskSuccess(t *testing.T) {
//	err := (&TaskAttributes{}).ValidatePluggableTask()
//	assert.Nil(t, err)
//
//}
