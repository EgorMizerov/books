package handlers

import (
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	logcontext "github.com/egormizerov/books/pkg/log/context"
)

func TestMiddlewareLoggerInContext(t *testing.T) {
	logger := logrus.New()
	request := &http.Request{}
	contextWithLogger := logcontext.WithLogger(request.Context(), logrus.NewEntry(logger))

	result := MiddlewareLoggerInContext(request, logger)

	assert.Equal(t, request.WithContext(contextWithLogger), result)
}
