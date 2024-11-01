package utils

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func IsProduction() bool {
	return viper.GetString("ENVIRONMENT") == "production"
}

func GetRegion() string {
	return viper.GetString("REGION")
}

func LogTimeTaken(functionName string, startTime time.Time) {
	logrus.Debugf("%s: %v", functionName, time.Since(startTime))
}
