package core

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/meorkamil/pgpool-consul/internal/consul"
	"github.com/meorkamil/pgpool-consul/internal/model"
	"github.com/meorkamil/pgpool-consul/internal/pgpool"
)

type core struct {
	config *model.Config
}

func NewPgpoolConsul(c *model.Config) *core {
	return &core{
		config: c,
	}
}

func (c *core) Run() error {
	slog.Info(fmt.Sprintf("Starting pgpool-consul: %s", c.config.Version))

	// Create a channel to listen for shutdown signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a context that can be canceled for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		// Create a channel for communication between goroutines
		pgpoolChan := make(chan string, 1)
		pgpool := pgpool.NewPgPool(*c.config)

		// Create consul client, TODO: remove it to diff goroutine
		cnsl, err := consul.NewConsul(*c.config)
		if err != nil {
			slog.Error(fmt.Sprintf("Consul client: %s", err))
		}

		// Start goroutine for pgpool
		go pgpool.Run(pgpoolChan)

		// Start goroutine for shutdown
		go c.handleShutdown(signalChan, cancel)

		select {
		// Handle the pgpool status received from the channel
		case pgpoolStat := <-pgpoolChan:
			// Create a consul client
			if err := cnsl.RegisterService(pgpoolStat); err != nil {
				slog.Error(fmt.Sprintf("Consul registration: %s", err))
			}

		// Gracefully shutdown and close resource
		case <-ctx.Done():
			// Deregister service
			if err := cnsl.RegisterService("SHUTDOWN"); err != nil {
				slog.Error(fmt.Sprintf("Consul registration: %s", err))
			}

			// Proceed with shutdown
			slog.Info("Shuting down")
			close(pgpoolChan)

			return nil
		}

		// Keep as a long runnign process
		time.Sleep(time.Duration(c.config.Global.Interval) * time.Second)
	}
}
