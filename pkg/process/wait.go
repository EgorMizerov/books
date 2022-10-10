package process

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForTermination() {
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGTERM)
	<-stopSignal
}
