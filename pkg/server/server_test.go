package server

import (
	"context"
	"errors"
	"github.com/egormizerov/books/app/handlers"
	"github.com/egormizerov/books/pkg/server/mocks"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

type ServerTests struct {
	suite.Suite
	server         Server
	httpServerMock *mocks.HttpServer
	testError      error
	context        context.Context
}

func TestServer(t *testing.T) {
	suite.Run(t, new(ServerTests))
}

func (self *ServerTests) SetupTest() {
	self.httpServerMock = mocks.NewHttpServer(self.T())
	self.server = Server{httpServer: self.httpServerMock}
	self.testError = errors.New("test_error")
	self.context = context.Background()
}

func (self *ServerTests) TestNewServer() {
	addr := "localhost:8080"
	handler := &handlers.Handler{}

	result := NewServer(addr, handler)

	self.Equal(&Server{
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       10 * time.Second,
		},
	}, result)
}

func (self *ServerTests) TestListen() {
	self.httpServerMock.
		On("ListenAndServe").
		Return(nil)

	err := self.server.Listen()

	self.NoError(err)
}

func (self *ServerTests) TestListenReturnsError() {
	self.httpServerMock.
		On("ListenAndServe").
		Return(self.testError)

	err := self.server.Listen()

	self.EqualError(err, self.testError.Error())
}

func (self *ServerTests) TestShutdown() {
	self.httpServerMock.
		On("Shutdown", self.context).
		Return(nil)

	err := self.server.Shutdown(self.context)

	self.NoError(err)
}

func (self *ServerTests) TestShutdownReturnsError() {
	self.httpServerMock.
		On("Shutdown", self.context).
		Return(self.testError)

	err := self.server.Shutdown(self.context)

	self.EqualError(err, self.testError.Error())
}
