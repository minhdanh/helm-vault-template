package main

import (
	"os"

	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

var Version = "0.0.1"

func createRenderer(c *cli.Context) (*renderer, error) {
	cfg := rendererConfig{
		vault: vaultConfig{
			endpoint: c.GlobalString("vault.address"),
			token:    c.GlobalString("vault.token"),
		},
	}

	return NewRenderer(cfg)
}

func NewYamlSourceFromFile(file string) func(context *cli.Context) (altsrc.InputSourceContext, error) {
	return func(context *cli.Context) (altsrc.InputSourceContext, error) {
		return altsrc.NewYamlSourceFromFile(file)
	}
}

func main() {
	app := cli.NewApp()
	app.Version = Version
	app.Usage = "Render template with values from Vault to a file"

	flags := []cli.Flag{
		cli.StringFlag{
			Name:   "vault.token",
			Usage:  "Vault API token",
			EnvVar: "VAULT_TOKEN",
		},
		altsrc.NewStringFlag(cli.StringFlag{
			Name:   "vault.address",
			Usage:  "Vault API endpoint",
			EnvVar: "VAULT_ADDR",
		}),
	}

	app.Flags = flags

	app.Commands = []cli.Command{
		{
			Name:  "render",
			Usage: "render a template with values from Vault into a file",
			Action: func(c *cli.Context) error {
				if c.NArg() < 2 {
					cli.ShowCommandHelpAndExit(c, "render", 1)
				}

				renderer, err := createRenderer(c)

				if err != nil {
					return err
				}

				inputFile := c.Args().Get(0)
				outputFile := c.Args().Get(1)

				return renderer.renderSingleFile(inputFile, outputFile)
			},
			ArgsUsage: "[input file] [output file]",
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		panic(err)
	}
}
