package logger

import "github.com/sirupsen/logrus"

func ConfigureLogger() {
	logrus.SetReportCaller(true)
}
