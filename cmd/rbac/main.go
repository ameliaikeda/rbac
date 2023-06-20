package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/ameliaikeda/rbac/generator"
	"github.com/ameliaikeda/rbac/generator/yaml"
)

func main() {
	app := &cli.App{
		Name:  "rbac",
		Usage: "generate roles and permissions for an application",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "rbac.yaml",
				Usage:   "yaml config file for generation",
			},
			&cli.StringFlag{
				Name:  "path",
				Value: "./rbac",
				Usage: "The path to the folder that should contain the role/permission packages.",
			},
		},
		Action: func(context *cli.Context) error {
			basePath, err := filepath.Abs(context.String("path"))
			if err != nil {
				return err
			}

			if err := generator.Run(context.Context,
				yaml.OptionsFromYAML(context.String("config")),
				generator.BasePath(basePath),
				generator.ResolvePaths,
			); err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
