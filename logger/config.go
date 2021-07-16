package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.Level(viper.GetInt("LOGGER.LOG_LEVEL")))
}
