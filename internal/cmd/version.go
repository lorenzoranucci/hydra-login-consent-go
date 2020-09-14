package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

func getVersionCommand(version string) cli.Command {
	return cli.Command{
		Name: "version",
		Action: func(_ *cli.Context) error {
			fmt.Println(version)
			return nil
		},
		Usage: "Print the version of the IdentityProvider binary",
	}
}
