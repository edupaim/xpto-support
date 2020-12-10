package app

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

var rootCmd = &cobra.Command{
	Run:   runApplication,
	Short: "Run XPTO support application",
	Long:  "",
}

func Execute() error {
	return rootCmd.Execute()
}

func runApplication(cmd *cobra.Command, args []string) {
	if err := initConfig(); err != nil {
		logrus.WithError(err).Fatalln("load config")
		return
	}
	api, err := initializeApi()
	if err != nil {
		logrus.WithError(err).Fatalln("initialize application")
		return
	}
	errChan := api.run()
	interruptAppChan := make(chan os.Signal)
	signal.Notify(interruptAppChan, os.Interrupt)
	defer api.shutdown()
	select {
	case err := <-errChan:
		logrus.Fatalln(err.Error())
	case sig := <-interruptAppChan:
		logrus.WithField("signal", sig.String()).Debugln("receive signal to interrupt application")
		logrus.Infoln("gracefully shutdown")
	}
	return
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./config.json", "config file (default is $HOME/config.json)")
}
