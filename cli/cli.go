package cli

import (
	"encoding/json"
	"fmt"

	"github.com/beamly/go-gocd/gocd"
	"github.com/urfave/cli"
)

// GetCliCommands returns a list of all CLI Command structs
func GetCliCommands() []cli.Command {
	return []cli.Command{
		*configureCommand(),
		*listAgentsCommand(),
		*listPipelineTemplatesCommand(),
		*getAgentCommand(),
		*getPipelineTemplateCommand(),
		*createPipelineTemplateCommand(),
		*updateAgentCommand(),
		*updateAgentsCommand(),
		*updatePipelineConfigCommand(),
		*updatePipelineTemplateCommand(),
		*deleteAgentCommand(),
		*deleteAgentsCommand(),
		*deletePipelineTemplateCommand(),
		*deletePipelineConfigCommand(),
		*listPipelineGroupsCommand(),
		*getPipelineHistoryCommand(),
		*getPipelineCommand(),
		*createPipelineConfigCommand(),
		*getPipelineStatusCommand(),
		*pausePipelineCommand(),
		*unpausePipelineCommand(),
		*releasePipelineLockCommand(),
		*getConfigurationCommand(),
		*encryptCommand(),
		*getVersionCommand(),
		*listPluginsCommand(),
		*getPluginCommand(),
		*listScheduledJobsCommand(),
		*getPipelineConfigCommand(),
		*listEnvironmentsCommand(),
		*getEnvironmentCommand(),
		*patchEnvironmentCommand(),
		*addPipelinesToEnvironmentCommand(),
		*removePipelinesFromEnvironmentCommand(),
		*listPropertiesCommand(),
		*createPropertyCommand(),
		*createRoleCommand(),
		*listRoleCommand(),
		*getRoleCommand(),
		*updateRoleCommand(),
		*deleteRoleCommand(),
	}
}

// NewCliClient creates a new gocd client for use by cli actions.
func NewCliClient(c *cli.Context) (cl *gocd.Client, err error) {
	var profile string

	if profile = c.String("profile"); profile == "" {
		profile = "default"
	}

	cfg := &gocd.Configuration{}
	err = gocd.LoadConfigByName(profile, cfg)

	setStringFromContext(&cfg.Server, "server", c)

	if cfg.Server == "" {
		if err == nil {
			// If we didn't have any errors, and our server is empty, use the local.
			cfg.Server = "https://127.0.0.1:8154/go/"
		} else {
			return
		}
		return
	}

	setStringFromContext(&cfg.Username, "username", c)
	setStringFromContext(&cfg.Password, "password", c)

	cfg.SkipSslCheck = cfg.SkipSslCheck || c.Bool("skip_ssl_check")

	return cfg.Client(), nil
}

func setStringFromContext(dest *string, key string, c *cli.Context) {
	var value string
	if value = c.String(key); value != "" {
		*dest = value
	}
}

func handleOutput(r interface{}, reqType string) cli.ExitCoder {
	o := map[string]interface{}{
		fmt.Sprintf("%s-response", reqType): r,
	}
	b, err := json.MarshalIndent(o, "", "    ")
	if err != nil {
		return NewCliError(reqType, nil, err)
	}

	fmt.Println(string(b))
	return nil
}

// ActionWrapperFunc describes the callback provided to ActionWrapper
type ActionWrapperFunc func(client *gocd.Client, c *cli.Context) (interface{}, *gocd.APIResponse, error)

// ActionWrapper handles the deferencing, and casting of the client object, and error handling.
func ActionWrapper(callback ActionWrapperFunc) interface{} {
	return func(c *cli.Context) error {
		cl := c.App.Metadata["c"].(func(c *cli.Context) (*gocd.Client, error))
		client, err := cl(c.Parent())
		if err != nil {
			return NewCliError(c.Command.Name, nil, err)
		}
		v, resp, err := callback(client, c)
		if err != nil {
			return NewCliError(c.Command.Name, resp, err)
		}
		return handleOutput(v, c.Command.Name)
	}
}
