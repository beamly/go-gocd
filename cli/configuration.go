package cli

import (
	"context"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/urfave/cli"
)

// List of command name and descriptions
const (
	GetConfigurationCommandName  = "get-configuration"
	GetConfigurationCommandUsage = "Get GoCD server configuration. This is the cruise-config.xml file. It is exposed in a json format to enable a consistent format. This API is for read-only purposes and not intended as an interface to modify the config."
	GetVersionCommandName        = "get-version"
	GetVersionCommandUsage       = "Return the Version for the GoCD Server"
)

// GetConfigurationAction gets a list of agents and return them.
func getConfigurationAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	return client.Configuration.Get(context.Background())
}

// GetVersionAction returns version information about GoCD
func getVersionAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	return client.Configuration.GetVersion(context.Background())
}

// GetConfigurationCommand handles the interaction between the cli flags and the action handler for delete-agents
func getConfigurationCommand() *cli.Command {
	return &cli.Command{
		Name:     GetConfigurationCommandName,
		Usage:    GetConfigurationCommandUsage,
		Action:   ActionWrapper(getConfigurationAction),
		Category: "Configuration",
	}
}

// GetVersionCommand handles the interaction between the cli flags and the action handler for delete-agents
func getVersionCommand() *cli.Command {
	return &cli.Command{
		Name:     GetVersionCommandName,
		Usage:    GetVersionCommandUsage,
		Action:   ActionWrapper(getVersionAction),
		Category: "Configuration",
	}
}
