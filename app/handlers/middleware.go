package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"

	logcontext "github.com/egormizerov/books/pkg/log/context"
)

func MiddlewareLoggerInContext(request *http.Request, logger *logrus.Logger) *http.Request {
	contextWithLogger := logcontext.WithLogger(request.Context(), logrus.NewEntry(logger))
	return request.WithContext(contextWithLogger)
}
