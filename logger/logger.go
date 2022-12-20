package logger

import (
        "io"
        logging "log"
        "os"

        "github.com/sirupsen/logrus"
)

var (
        log *logrus.Logger
)

func init() {
        f, err := os.OpenFile("logs/application.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
        if err != nil {
                logging.Fatalf("error opening file: %v", err)
        }
        log = logrus.New()
        log.Formatter = &logrus.TextFormatter{}
        // log.SetReportCaller(true)
        mw := io.MultiWriter(os.Stdout, f)
        log.SetOutput(mw)
}

// Info ...
func Info(format string, v ...interface{}) {
        log.Infof(format, v...)
}

// Warn ...
func Warn(format string, v ...interface{}) {
        log.Warnf(format, v...)
}

// Error ...
func Error(format string, v ...interface{}) {
        log.Errorf(format, v...)
}

// Fatal ...
func Fatal(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

var (

        // RecommendationError ...
        RecommendationError = "%v type=recommendation.error"

        // RecommendationWarn ...
        RecommendationWarn = "%v type=recommendation.warn"

        // RecommendationInfo ...
        RecommendationInfo = "%v type=recommendation.info"

	// RecommendationFatal ...
        RecommendationFatal = "%v type=recommendation.fatal"

        // ProvisionerWarn ...
        ProvisionerWarn = "%v type=provisioner.warn"

        // ProvisionerInfo ...
        ProvisionerInfo = "%v type=provisioner.info"

        // ProvisionerError ...
        ProvisionerError = "%v type=provisioner.error"

	// ProvisionerFatal ...
        ProvisionerFatal = "%v type=provisioner.fatal"
)

