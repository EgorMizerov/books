package log

import (
	"testing"

	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	pkgtesting "github.com/egormizerov/books/pkg/testing"
)

type LogrusTests struct {
	suite.Suite
}

func TestLogrus(t *testing.T) {
	suite.Run(t, new(LogrusTests))
}

func (self *LogrusTests) TestNewLogrusLoggerReturnsDefault() {
	logger := NewLogrusLogger(false, logrus.InfoLevel)

	self.LoggersEqual(logrus.New(), logger)
}

func (self *LogrusTests) TestNewLogrusLoggerWithJoonixFormatter() {
	expectedLogger := logrus.New()
	expectedFormatter := joonix.NewFormatter()
	expectedFormatter.TimestampFormat = logrusTimestampFormat
	expectedLogger.Formatter = expectedFormatter

	logger := NewLogrusLogger(true, logrus.InfoLevel)

	self.JoonixLoggersEqual(expectedLogger, logger)
}

func (self *LogrusTests) TestNewLogrusLoggerWithDebugLevel() {
	expectedLogger := logrus.New()
	expectedLogger.Level = logrus.DebugLevel

	logger := NewLogrusLogger(false, logrus.DebugLevel)

	self.LoggersEqual(expectedLogger, logger)
}

func (self *LogrusTests) LoggersEqual(expectedLogger, actualLogger *logrus.Logger) {
	expectedCopy, actualCopy := *expectedLogger, *actualLogger
	expectedCopy.ExitFunc, actualCopy.ExitFunc = nil, nil

	// self.Equal does not support comparison of function fields other than the value "nil".
	// So we check them separately above and set "nil" to safely use the Equal method.
	pkgtesting.FuncsEqual(self.T(), expectedLogger.ExitFunc, actualLogger.ExitFunc)
	self.Equal(expectedCopy, actualCopy)
}

func (self *LogrusTests) JoonixLoggersEqual(expectedLogger, actualLogger *logrus.Logger) {
	expectedCopy, actualCopy := *expectedLogger, *actualLogger
	expectedCopy.Formatter, actualCopy.Formatter = nil, nil

	// self.Equal does not support comparison of function fields other than the value "nil".
	// So we check them separately above and set "nil" to safely use the Equal method.
	self.FormattersEqual(expectedLogger.Formatter, actualLogger.Formatter)
	self.LoggersEqual(&expectedCopy, &actualCopy)
}

func (self *LogrusTests) FormattersEqual(expectedFormatter, actualFormatter logrus.Formatter) {
	expected, actual := expectedFormatter.(*joonix.Formatter), actualFormatter.(*joonix.Formatter)
	expectedCopy, actualCopy := *expected, *actual
	expectedCopy.TimestampFormat, actualCopy.TimestampFormat = nil, nil

	// self.Equal does not support comparison of function fields other than the value "nil".
	// So we check them separately above and set "nil" to safely use the Equal method.
	pkgtesting.FuncsEqual(self.T(), expected.TimestampFormat, actual.TimestampFormat)
	self.Equal(expectedCopy, actualCopy)
}
