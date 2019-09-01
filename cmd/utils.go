package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func waitQuit() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		_ = <-sigs
		fmt.Println()
		done <- true
	}()

	<-done
}
