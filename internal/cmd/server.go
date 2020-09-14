package cmd

import (
	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/infrastructure/http"
	"github.com/urfave/cli"
)

func getServerCommand(baseFlags []cli.Flag, baseBeforeFunc cli.BeforeFunc) cli.Command {
	return cli.Command{
		Name:   "server",
		Action: runServer,
		Usage:  "Run the http server which expose IdentityProvider API",
		Flags: append(
			baseFlags,
			cli.IntFlag{
				Name:   "port",
				Value:  9020,
				Usage:  "Server port",
				EnvVar: "PPRO_IDENTITY_PROVIDER_PORT",
			},
			cli.IntFlag{
				Name:   "healthcheck-port",
				Value:  8086,
				Usage:  "Healthcheck Server port",
				EnvVar: "PPRO_IDENTITY_PROVIDER_HEALTHCHECK_PORT",
			},
		),
		Before: baseBeforeFunc,
	}
}

func runServer(c *cli.Context) error {
	//go healthcheck.NewServer(c.Int("healthcheck-port"), app.GetServiceLocator().HydraAdmin()).Run()

	http.NewServer(
		c.Int("port"),
	).Run()

	return nil
}
