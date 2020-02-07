package gocd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStages(t *testing.T) {
	t.Run("Validate", testStageValidate)
	t.Run("JSONStringFail", testStageJSONStringFail)
	t.Run("JSONString", testStageJSONString)
}

func testStageJSONStringFail(t *testing.T) {
	s := Stage{Approval: &Approval{Type: "success"}}
	_, err := s.JSONString()
	assert.EqualError(t, err, "`gocd.Stage.Name` is empty")
}

func testStageJSONString(t *testing.T) {
	s := Stage{
		Name:     "test-stage",
		Approval: &Approval{Type: "success"},
		Jobs:     []*Job{{Name: "test-job"}},
	}
	j, err := s.JSONString()
	if err != nil {
		assert.Nil(t, err)
	}
	assert.Equal(
		t, `{
  "name": "test-stage",
  "fetch_materials": false,
  "clean_working_directory": false,
  "never_cleanup_artifacts": false,
  "approval": {
    "type": "success",
    "authorization": {
      "users": [],
      "roles": []
    }
  },
  "jobs": [
    {
      "name": "test-job"
    }
  ]
}`, j)
}

func testStageValidate(t *testing.T) {
	s := Stage{}

	err := s.Validate()
	assert.EqualError(t, err, "`gocd.Stage.Name` is empty")

	s.Name = "test-stage"
	err = s.Validate()
	assert.EqualError(t, err, "At least one `Job` must be specified")

	s.Jobs = []*Job{{}}
	err = s.Validate()
	assert.EqualError(t, err, "`gocd.Jobs.Name` is empty")

	s.Jobs[0].Name = "test-job"
	err = s.Validate()
	assert.Nil(t, err)
}
