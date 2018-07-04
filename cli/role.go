package cli

import (
	"context"
	"github.com/beamly/go-gocd/gocd"
	"github.com/urfave/cli"
	"errors"
	"io/ioutil"
	"encoding/json"
)

// List of command name and descriptions
const (
	CreateRoleCommandName  = "create-role"
	CreateRoleCommandUsage = "Create a role"
	ListRoleCommandName    = "list-role"
	ListRoleCommandUsage   = "List all the roles"
	GetRoleCommandName     = "get-role"
	GetRoleCommandUsage    = "Get a Role"
	DeleteRoleCommandName  = "delete-role"
	DeleteRoleCommandUsage = "Delete a role"
	UpdateRoleCommandName  = "update-role"
	UpdateRoleCommandUsage = "Update a Role"
)

func createRoleAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	name := c.String("name")
	if name == "" {
		return nil, nil, NewFlagError("name")
	}

	roleJson := c.String("role-json")
	roleFile := c.String("role-file")
	if roleJson == "" && roleFile == "" {
		return nil, nil, errors.New("One of '--role-file' or '--role-json' must be specified")
	}

	if roleJson != "" && roleFile != "" {
		return nil, nil, errors.New("Only one of '--role-file' or '--role-json' can be specified")
	}

	var rf []byte
	if roleFile != "" {
		rf, err = ioutil.ReadFile(roleFile)
		if err != nil {
			return nil, nil, err
		}
	} else {
		rf = []byte(roleJson)
	}
	role := &gocd.Role{}
	err = json.Unmarshal(rf, &role)
	if err != nil {
		return nil, nil, err
	}

	role.Name = name

	return client.Roles.Create(context.Background(), role)

}

func getRoleAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	name := c.String("name")
	if name == "" {
		return nil, nil, NewFlagError("name")
	}

	return client.Role.Get(context.Background(), c.String("name"))
}

// ListRoleAction retrieves all role configurations
func listRoleAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	return client.Role.List(context.Background())
}

func createRoleCommand() *cli.Command {
	return &cli.Command{
		Name:     CreateRoleCommandName,
		Usage:    CreateRoleCommandUsage,
		Category: "Role",
		Action:   ActionWrapper(createRoleAction),
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
			cli.StringFlag{Name: "role-json", Usage: "A JSON string describing the role configuration"},
			cli.StringFlag{Name: "role-file", Usage: "Path to a JSON file describing the role configuration"},
		},
	}
}

func listRoleCommand() *cli.Command {
	return &cli.Command{
		Name:     ListRoleCommandName,
		Usage:    ListRoleCommandUsage,
		Category: "Role",
		Action:   ActionWrapper(listRoleAction),
	}
}

func getRoleCommand() *cli.Command {
	return &cli.Command{
		Name:     GetRoleCommandName,
		Usage:    GetRoleCommandUsage,
		Category: "Role",
		Action:   ActionWrapper(getRoleAction),
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
	}
}

func deleteRoleCommand() *cli.Command {
	return &cli.Command{
		Name:     DeleteRoleCommandName,
		Usage:    DeleteRoleCommandUsage,
		Category: "Role",
		Action:   ActionWrapper(deleteRoleAction),
	}
}

func updateRoleCommand() *cli.Command {
	return &cli.Command{
		Name:     UpdateRoleCommandName,
		Usage:    UpdateRoleCommandUsage,
		Category: "Role",
		Action:   ActionWrapper(updateRoleAction),
	}
}
