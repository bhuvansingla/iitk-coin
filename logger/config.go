package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ConfigureLogger() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.Level(viper.GetInt("LOGGER.LOG_LEVEL")))
}
