package cli

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/beamly/go-gocd/gocd"
	"github.com/urfave/cli"
	"io/ioutil"
)

// List of command name and descriptions
const (
	CreateRoleCommandName  = "create-role"
	CreateRoleCommandUsage = "Create a role"
	ListRoleCommandName    = "list-roles"
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

	roleJSON := c.String("role-json")
	roleFile := c.String("role-file")
	if roleJSON == "" && roleFile == "" {
		return nil, nil, errors.New("One of '--role-file' or '--role-json' must be specified")
	}

	if roleJSON != "" && roleFile != "" {
		return nil, nil, errors.New("Only one of '--role-file' or '--role-json' can be specified")
	}

	var rf []byte
	if roleFile != "" {
		rf, err = ioutil.ReadFile(roleFile)
		if err != nil {
			return nil, nil, err
		}
	} else {
		rf = []byte(roleJSON)
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

	getResponse, resp, err := client.Roles.Get(context.Background(), name)
	if resp.HTTP.StatusCode != 404 {
		getResponse.RemoveLinks()
	}
	return getResponse, resp, err
}

// ListRoleAction retrieves all role configurations
func listRoleAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	return client.Roles.List(context.Background())
}

// UpdateRoleAction handles the interaction between the cli flags and the action handler for
// update-role-action
func updateRoleAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	var name, version string

	if name = c.String("name"); name == "" {
		return nil, nil, NewFlagError("name")
	}

	if version = c.String("role-version"); version == "" {
		return nil, nil, NewFlagError("role-version")
	}

	role := c.String("role")
	roleFile := c.String("role-file")
	if role == "" && roleFile == "" {
		return nil, nil, errors.New("One of '--role-file' or '--role' must be specified")
	}

	if role != "" && roleFile != "" {
		return nil, nil, errors.New("Only one of '--role-file' or '--role' can be specified")
	}

	var pf []byte
	if roleFile != "" {
		pf, err = ioutil.ReadFile(roleFile)
		if err != nil {
			return nil, nil, err
		}
	} else {
		pf = []byte(role)
	}
	p := &gocd.Role{
		Version: version,
	}
	err = json.Unmarshal(pf, &p)
	if err != nil {
		return nil, nil, err
	}

	return client.Roles.Update(context.Background(), name, p)
}

func deleteRoleAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	name := c.String("name")
	if name == "" {
		return nil, nil, NewFlagError("name")
	}

	deleteResponse, resp, err := client.Roles.Delete(context.Background(), name)
	if resp.HTTP.StatusCode == 406 {
		err = errors.New(deleteResponse)
	}
	return deleteResponse, resp, err
}
func createRoleCommand() *cli.Command {
	return &cli.Command{
		Name:     CreateRoleCommandName,
		Usage:    CreateRoleCommandUsage,
		Category: "Roles",
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
		Category: "Roles",
		Action:   ActionWrapper(listRoleAction),
	}
}

func getRoleCommand() *cli.Command {
	return &cli.Command{
		Name:     GetRoleCommandName,
		Usage:    GetRoleCommandUsage,
		Category: "Roles",
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
		Category: "Roles",
		Action:   ActionWrapper(deleteRoleAction),
	}
}

func updateRoleCommand() *cli.Command {
	return &cli.Command{
		Name:     UpdateRoleCommandName,
		Usage:    UpdateRoleCommandUsage,
		Category: "Roles",
		Action:   ActionWrapper(updateRoleAction),
	}
}
