package cmd

import (
	"fmt"
	"os"

	"github.com/hegdeshashank73/glamr-backend/utils"
	"github.com/hegdeshashank73/glamr-backend/vendors"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "mono",
	Short: "CLI for running and managing Duggup Mono",
}
var ENVIRONMENT = os.Getenv("ENVIRONMENT")

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLogging)
	cobra.OnInitialize(vendors.Setup)
}

func Execute() error {

	return rootCmd.Execute()
}

func initConfig() {
	configName := "config"
	if os.Getenv("REGION") != "" {
		configName += "-"
		configName += os.Getenv("REGION")
	}

	viper.SetConfigName(configName) // name of your config file
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config") // path to your config file
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	viper.Set("REGION", os.Getenv("REGION"))
	os.Setenv("OTEL_EXPORTER_OTLP_HEADERS", viper.GetString("OTEL_EXPORTER_OTLP_HEADERS"))

	basepath := "."
	if utils.IsProduction() {
		basepath = "/home/ubuntu/mono"
	}
	viper.Set("BASEPATH", basepath)
}

func initLogging() {
	switch viper.GetString("LOG_LEVEL") {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	}

	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)
}
