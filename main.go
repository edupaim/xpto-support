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

const configFlag = "config"
const portFlag = "port"

func main() {
	cliApp := &cli.App{
		Name:   "xtpo support application",
		Usage:  "request drain for xpto application",
		Action: runApplication,
	}
	cliApp.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        configFlag + ", c",
			Usage:       "configuration file",
			Value:       "./config.json",
			Destination: &fileName,
		},
		&cli.IntFlag{
			Name:        portFlag + ", p",
			Usage:       "application port",
			Value:       5051,
			Destination: &port,
		},
	}
	err := cliApp.Run(os.Args)
	if err != nil {
		logrus.Fatalln(err)
	}
}

func runApplication(c *cli.Context) error {
	config, err := app.InitConfig(c.String("config"))
	if err != nil {
		return err
	}
	if port := c.Int(portFlag); port != 0 {
		config.WithPort(port)
	}
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
