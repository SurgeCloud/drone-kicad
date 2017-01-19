package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {

	app := cli.NewApp()
	app.Name = "kicad plugin"
	app.Usage = "kicad plugin"
	app.Action = run
	app.Version = fmt.Sprintf("0.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "client.code",
			Usage:  "enterprise client code",
			EnvVar: "PLUGIN_CLIENT_CODE",
		},
		cli.StringFlag{
			Name:   "client.name",
			Usage:  "client name",
			EnvVar: "PLUGIN_CLIENT_NAME",
		},
		cli.StringFlag{
			Name:   "project.code",
			Usage:  "enterprise project code",
			EnvVar: "PLUGIN_PROJECT_CODE",
		},
		cli.StringFlag{
			Name:   "project.name",
			Usage:  "project name",
			EnvVar: "PLUGIN_PROJECT_NAME",
		},
		cli.BoolFlag{
			Name:   "options.schematic",
			Usage:  "generate schematic",
			EnvVar: "PLUGIN_SCHEMATIC",
		},
		cli.BoolFlag{
			Name:   "options.bom",
			Usage:  "generate bom",
			EnvVar: "PLUGIN_BOM",
		},
		cli.StringFlag{
			Name:   "options.gerber",
			Usage:  "generate gerber files",
			EnvVar: "PLUGIN_GERBER",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {

	plugin := Plugin{
		Client: Client{
			Code: c.String("client.code"),
			Name: c.String("client.name"),
		},
		Project: Project{
			Code: c.String("project.code"),
			Name: c.String("project.name"),
		},
		Options: Options{
			Sch: c.Bool("options.schematic"),
			Bom: c.Bool("options.bom"),
			//Grb:    nil,
			GrbGen: c.IsSet("options.gerber"),
		},
	}

	if plugin.Options.GrbGen {
		err := json.Unmarshal([]byte(c.String("options.gerber")), &plugin.Options.Grb)
		if err != nil {
			return err
		}
	}

	return plugin.Exec()
}
