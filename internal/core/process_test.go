package core

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestHandleShutdown(t *testing.T) {
	// Create signal channel for testing
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Create context for testing
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init empty core
	c := &core{}

	// Run handle shutdown
	go c.handleShutdown(signalChan, cancel)

	// Send signal into the channel
	signalChan <- syscall.SIGINT

	// Delay to let the proceess end
	time.Sleep(100 * time.Millisecond)
}
