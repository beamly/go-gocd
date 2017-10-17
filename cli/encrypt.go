package cli

import (
	"context"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/urfave/cli"
)

// List of command name and descriptions
const (
	EncryptCommandName  = "encrypt"
	EncryptCommandUsage = "Encrypt a value for use in GoCD configurations"
)

// EncryptAction gets a list of agents and return them.
func encryptAction(client *gocd.Client, c *cli.Context) (r interface{}, resp *gocd.APIResponse, err error) {
	var value string
	if value = c.String("value"); value == "" {
		return nil, nil, NewFlagError("value")
	}

	return client.Encryption.Encrypt(context.Background(), value)
}

// EncryptCommand checks a template-name is provided and that the response is a 2xx response.
func encryptCommand() *cli.Command {
	return &cli.Command{
		Name:     EncryptCommandName,
		Usage:    EncryptCommandUsage,
		Action:   ActionWrapper(encryptAction),
		Category: "Encryption",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "value"},
		},
	}
}
