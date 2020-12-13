package main

import (
	"edupaim/xpto-support/app"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
)

var fileName string
var port int
var debug bool

const configFlag = "config"
const portFlag = "port"
const debugFlag = "debug"

func main() {
	cliApp := &cli.App{
		Name:  "xtpo support application",
		Usage: "request drain for xpto application",
		Action: func(c *cli.Context) error {
			return runApplication()
		},
	}
	cliApp.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        configFlag + ", c",
			Usage:       "configuration file",
			Value:       "./config.json",
			Required:    true,
			Destination: &fileName,
		},
		&cli.IntFlag{
			Name:        portFlag + ", p",
			Usage:       "application port",
			Value:       0,
			Required:    true,
			Destination: &port,
		},
		&cli.BoolFlag{
			Name:        debugFlag + ", d",
			Usage:       "debug flag",
			Value:       false,
			Destination: &debug,
		},
	}
	err := cliApp.Run(os.Args)
	if err != nil {
		logrus.Fatalln(err)
	}
}

func runApplication() error {
	config, err := app.InitConfig(fileName, debug)
	if err != nil {
		return err
	}
	config.WithPort(port)
	api, err := app.InitializeApi(config)
	if err != nil {
		return err
	}
	errChan := api.Run()
	return waitForGracefullyShutdown(api, errChan)
}

func waitForGracefullyShutdown(api *app.Api, errChan <-chan error) error {
	interruptAppChan := make(chan os.Signal)
	signal.Notify(interruptAppChan, os.Interrupt)
	defer api.Shutdown()
	select {
	case err := <-errChan:
		return err
	case <-interruptAppChan:
	}
	return nil
}
