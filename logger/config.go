package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	// logrus.SetReportCaller(true)

	logrus.SetLevel(logrus.Level(viper.GetInt64("LOGGER.LOG_LEVEL")))

	f, err := os.OpenFile("iitkcoin.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	logrus.SetOutput(f)

}
