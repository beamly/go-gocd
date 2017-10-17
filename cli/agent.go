package cli

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/urfave/cli"
)

// List of command name and descriptions
const (
	ListAgentsCommandName    = "list-agents"
	ListAgentsCommandUsage   = "List GoCD build agents."
	GetAgentCommandName      = "get-agent"
	GetAgentCommandUsage     = "Get Agent by UUID"
	UpdateAgentCommandName   = "update-agent"
	UpdateAgentCommandUsage  = "Update Agent"
	DeleteAgentCommandName   = "delete-agent"
	DeleteAgentCommandUsage  = "Delete Agent"
	UpdateAgentsCommandName  = "update-agents"
	UpdateAgentsCommandUsage = "Bulk Update Agents"
	DeleteAgentsCommandName  = "delete-agents"
	DeleteAgentsCommandUsage = "Bulk Delete Agents"
	agentCategory            = "Agents"
)

// ListAgentsAction gets a list of agents and return them.
func listAgentsAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	agents, resp, err := client.Agents.List(context.Background())
	if err == nil {
		for _, agent := range agents {
			agent.RemoveLinks()
		}
	}
	return agents, resp, err
}

// GetAgentAction retrieves a single agent object.
func getAgentAction(client *gocd.Client, c *cli.Context) (v interface{}, resp *gocd.APIResponse, err error) {
	agent, r, err := client.Agents.Get(context.Background(), c.String("uuid"))
	if r.HTTP.StatusCode != 404 {
		agent.RemoveLinks()
	}
	return agent, r, err
}

// UpdateAgentAction updates a single agent.
func updateAgentAction(client *gocd.Client, c *cli.Context) (v interface{}, resp *gocd.APIResponse, err error) {

	if c.String("uuid") == "" {
		return nil, nil, NewFlagError("uuid")
	}

	if c.String("config") == "" {
		return nil, nil, NewFlagError("config")
	}

	a := &gocd.Agent{}
	b := []byte(c.String("config"))
	if err := json.Unmarshal(b, &a); err != nil {
		return nil, nil, err
	}

	return client.Agents.Update(context.Background(), c.String("uuid"), a)
}

// DeleteAgentAction delets an agent. Note: The agent must be disabled.
func deleteAgentAction(client *gocd.Client, c *cli.Context) (v interface{}, resp *gocd.APIResponse, err error) {
	if c.String("uuid") == "" {
		return nil, nil, NewFlagError("uuid")
	}

	deleteResponse, r, err := client.Agents.Delete(context.Background(), c.String("uuid"))
	if r.HTTP.StatusCode == 406 {
		err = errors.New(deleteResponse)
	}
	return deleteResponse, r, err
}

// UpdateAgentsAction updates a single agent.
func updateAgentsAction(client *gocd.Client, c *cli.Context) (v interface{}, resp *gocd.APIResponse, err error) {

	u := gocd.AgentBulkUpdate{}
	if o := c.String("operations"); o != "" {
		b := []byte(o)
		op := gocd.AgentBulkOperationsUpdate{}
		if err := json.Unmarshal(b, &op); err != nil {
			return nil, nil, err
		}
		u.Operations = &op
	}

	var uuids []string
	if uuids = c.StringSlice("uuid"); len(uuids) == 0 {
		return nil, nil, NewFlagError("uuid")
	}
	u.Uuids = uuids

	if state := c.String("state"); state != "" {
		u.AgentConfigState = c.String("state")
	}

	updateResponse, r, err := client.Agents.BulkUpdate(context.Background(), u)
	if r.HTTP.StatusCode == 406 {
		err = errors.New(updateResponse)
	}
	return updateResponse, r, err
}

// DeleteAgentsAction must be implemented.
func deleteAgentsAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	return nil, nil, nil
}

// ListAgentsCommand checks a template-name is provided and that the response is a 2xx response.
func listAgentsCommand() *cli.Command {
	return &cli.Command{
		Name:     ListAgentsCommandName,
		Usage:    ListAgentsCommandUsage,
		Action:   ActionWrapper(listAgentsAction),
		Category: agentCategory,
	}
}

// GetAgentCommand handles the interaction between the cli flags and the action handler for get-agent
func getAgentCommand() *cli.Command {
	return &cli.Command{
		Name:     GetAgentCommandName,
		Usage:    GetAgentCommandUsage,
		Action:   ActionWrapper(getAgentAction),
		Category: agentCategory,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "uuid, u", Usage: "GoCD Agent UUID"},
		},
	}
}

// UpdateAgentCommand handles the interaction between the cli flags and the action handler for update-agent
func updateAgentCommand() *cli.Command {
	return &cli.Command{
		Name:     UpdateAgentCommandName,
		Usage:    UpdateAgentCommandUsage,
		Action:   ActionWrapper(updateAgentAction),
		Category: agentCategory,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "uuid, u", Usage: "GoCD Agent UUID"},
			cli.StringFlag{Name: "config, c", Usage: "JSON encoded config for agent update."},
		},
	}
}

// DeleteAgentCommand handles the interaction between the cli flags and the action handler for delete-agent
func deleteAgentCommand() *cli.Command {
	return &cli.Command{
		Name:     DeleteAgentCommandName,
		Usage:    DeleteAgentCommandUsage,
		Action:   ActionWrapper(deleteAgentAction),
		Category: agentCategory,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "uuid, u", Usage: "GoCD Agent UUID"},
		},
	}
}

// UpdateAgentsCommand handles the interaction between the cli flags and the action handler for update-agents
func updateAgentsCommand() *cli.Command {
	return &cli.Command{
		Name:     UpdateAgentsCommandName,
		Usage:    UpdateAgentsCommandUsage,
		Action:   ActionWrapper(updateAgentsAction),
		Category: agentCategory,
		Flags: []cli.Flag{
			cli.StringSliceFlag{Name: "uuid", Usage: "GoCD Agent UUIDs"},
			cli.StringFlag{Name: "state", Usage: "Whether agents are enabled or disabled. Allowed values 'Enabled','Disabled'."},
			cli.StringFlag{Name: "operations", Usage: "JSON encoded config for bulk operation updates."},
		},
	}
}

// DeleteAgentsCommand handles the interaction between the cli flags and the action handler for delete-agents
func deleteAgentsCommand() *cli.Command {
	return &cli.Command{
		Name:     DeleteAgentsCommandName,
		Usage:    DeleteAgentsCommandUsage,
		Action:   ActionWrapper(deleteAgentsAction),
		Category: agentCategory,
		Flags: []cli.Flag{
			cli.StringSliceFlag{Name: "uuid", Usage: "GoCD Agent UUIDs"},
		},
	}
}
