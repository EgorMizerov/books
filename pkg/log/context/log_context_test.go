package context

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type LogContextTests struct {
	suite.Suite
}

func TestLogContext(t *testing.T) {
	suite.Run(t, new(LogContextTests))
}

func (self *LogContextTests) TestFromContext() {
	logger := logrus.NewEntry(logrus.StandardLogger())
	contextWithLogger := context.WithValue(context.Background(), loggerKey, logger)

	self.Equal(
		FromContext(contextWithLogger),
		logger,
	)
}

func (self *LogContextTests) TestFromContextWithoutLogger() {
	contextWithoutLogger := context.Background()

	self.Panics(func() {
		FromContext(contextWithoutLogger)
	})
}

func (self *LogContextTests) TestWithLogger() {
	logger := logrus.NewEntry(logrus.StandardLogger())
	emptyContext := context.Background()

	self.Equal(
		WithLogger(emptyContext, logger).Value(loggerKey),
		logger,
	)
}
