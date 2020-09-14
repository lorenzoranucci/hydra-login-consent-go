package cmd

import "github.com/urfave/cli"

func GetApp(version string) *cli.App {
	app := cli.NewApp()

	app.Version = version

	app.Name = "Identity Provider"
	app.Usage = ""

	app.HideVersion = true

	app.Commands = []cli.Command{
		getVersionCommand(app.Version),
		getServerCommand(app.Flags, app.Before),
	}

	return app
}

