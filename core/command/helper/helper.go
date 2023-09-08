package helper

import (
	"errors"
	"fmt"
	"github.com/Aleksao998/LightingUserVault/core/command"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ClientCloseResult struct {
	Message string `json:"message"`
}

func (r *ClientCloseResult) GetOutput() string {
	return r.Message
}

// HandleSignals is a helper method for handling signals sent to the console
// Like stop, error, etc.
func HandleSignals(
	closeFn func() error,
	outputter command.OutputFormatter,
) error {
	signalCh := getTerminationSignalCh()
	sig := <-signalCh

	closeMessage := fmt.Sprintf("\n[SIGNAL] Caught signal: %v\n", sig)
	closeMessage += "Gracefully shutting down client...\n"

	outputter.SetCommandResult(
		&ClientCloseResult{
			Message: closeMessage,
		},
	)
	outputter.WriteOutput()

	// Call the server close callback
	gracefulCh := make(chan struct{})

	go func() {
		if closeFn != nil {
			closeFn()
		}

		close(gracefulCh)
	}()

	select {
	case <-signalCh:
		return errors.New("shutdown by signal channel")
	case <-time.After(5 * time.Second):
		return errors.New("shutdown by timeout")
	case <-gracefulCh:
		return nil
	}
}

// GetTerminationSignalCh returns a channel to emit signals by ctrl + c
func getTerminationSignalCh() <-chan os.Signal {
	// wait for the user to quit with ctrl-c
	signalCh := make(chan os.Signal, 1)
	signal.Notify(
		signalCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	return signalCh
}
