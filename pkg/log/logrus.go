package log

import (
	"time"

	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
)

func NewLogrusLogger(enableJson bool, logLevel logrus.Level) *logrus.Logger {
	logger := logrus.New()
	if enableJson {
		formatter := joonix.NewFormatter()
		// Update TimestampFormat to be human readable.
		formatter.TimestampFormat = logrusTimestampFormat
		logger.Formatter = formatter
	}
	logger.SetLevel(logLevel)
	return logger
}

var logrusTimestampFormat = func(fields logrus.Fields, now time.Time) error {
	fields["time"] = now.Format(time.RFC3339)
	return nil
}
