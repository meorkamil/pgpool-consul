package core

import (
	"context"
	"log/slog"
	"os"
)

// Handler shutdown. TODO; will create a handler for process to reload the process
func (c *core) handleShutdown(signalChan chan os.Signal, cancel context.CancelFunc) {
	<-signalChan
	slog.Info("Received shutdown signal")
	cancel()
}
