package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"testing"
)

func TestRole(t *testing.T) {
	for _, envCmd := range []cli.Command{
		*createRoleCommand(),
		*listRoleCommand(),
		*getRoleCommand(),
		*deleteRoleCommand(),
		*updateRoleCommand(),
	} {
		assert.Equal(t, envCmd.Category, "Roles")
		assert.NotEmpty(t, envCmd.Name)
		assert.NotEmpty(t, envCmd.Usage)
	}
}
