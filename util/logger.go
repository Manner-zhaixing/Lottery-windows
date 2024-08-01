package util

import (
	"fmt"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var (
	LogRus *logrus.Logger
)

func InitLog(configFile string) {
	viper := CreateConfig(configFile)
	LogRus = logrus.New()
	// 根据配置文件中设置的日志等级字段，设置LogRus的日志等级
	switch strings.ToLower(viper.GetString("level")) {
	case "debug":
		LogRus.SetLevel(logrus.DebugLevel)
	case "info":
		LogRus.SetLevel(logrus.InfoLevel)
	case "warn":
		LogRus.SetLevel(logrus.WarnLevel)
	case "error":
		LogRus.SetLevel(logrus.ErrorLevel)
	case "panic":
		LogRus.SetLevel(logrus.PanicLevel)
	default:
		panic(fmt.Errorf("invalid log level %s", viper.GetString("level")))
	}

	LogRus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	logFile := ProjectRootPath + viper.GetString("file")
	fout, err := rotatelogs.New(
		logFile+".%Y%m%d%H",
		rotatelogs.WithLinkName(logFile),
		rotatelogs.WithRotationTime(1*time.Hour),
		rotatelogs.WithMaxAge(7*24*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	LogRus.SetOutput(fout)
	LogRus.SetReportCaller(true)
}
