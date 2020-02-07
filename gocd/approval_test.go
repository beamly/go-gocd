package gocd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApproval(t *testing.T) {
	approval := &Approval{
		Type: "success",
		Authorization: &Authorization{
			Roles: []string{"one"},
		},
	}
	assert.NotNil(t, approval.Authorization)
}
