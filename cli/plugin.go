package cli

import (
	"context"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/urfave/cli"
)

// List of command name and descriptions
const (
	ListPluginsCommandName  = "list-plugins"
	ListPluginsCommandUsage = "List all the Plugins"
	GetPluginCommandName    = "get-plugin"
	GetPluginCommandUsage   = "Get a Plugin"
)

// GetPluginAction retrieves a single plugin by name
func getPluginAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	name := c.String("name")
	if name == "" {
		return nil, nil, NewFlagError("name")
	}

	return client.Plugins.Get(context.Background(), c.String("name"))
}

// ListPluginsAction retrieves all plugin configurations
func listPluginsAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	return client.Plugins.List(context.Background())
}

// GetPluginCommand Describes the cli interface for the GetPluginAction
func getPluginCommand() *cli.Command {
	return &cli.Command{
		Name:     GetPluginCommandName,
		Usage:    GetPluginCommandUsage,
		Category: "Plugins",
		Action:   ActionWrapper(getPluginAction),
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
	}
}

// ListPluginsCommand Describes the cli interface for the ListPluginsCommand
func listPluginsCommand() *cli.Command {
	return &cli.Command{
		Name:     ListPluginsCommandName,
		Usage:    ListPluginsCommandUsage,
		Category: "Plugins",
		Action:   ActionWrapper(listPluginsAction),
	}
}
