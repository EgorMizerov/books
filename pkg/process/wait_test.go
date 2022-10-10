package process

import (
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/suite"
	"time"
)

type WaitTests struct {
	suite.Suite
}

func (self *WaitTests) TestProgramFinishedNormallyAfterSigterm() {
	go func() {
		time.Sleep(time.Millisecond)
		self.raise(syscall.SIGTERM)
	}()
	WaitForTermination()

	// If the test succeeded, the SIGTERM signal was intercepted.
}

func TestWait(t *testing.T) {
	suite.Run(t, new(WaitTests))
}

func (self *WaitTests) raise(sig os.Signal) {
	process, err := os.FindProcess(os.Getpid())
	self.Nil(err)
	self.Nil(process.Signal(sig))
}
