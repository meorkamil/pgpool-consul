package core

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/meorkamil/pgpool-consul/internal/consul"
	"github.com/meorkamil/pgpool-consul/internal/util"
)

const (
	TESTCONFIG = "../../config/config.yml"
)

func TestCore(t *testing.T) {
	// Init configuration file
	c, err := util.ConfigInit(TESTCONFIG, "test")
	if err != nil {
		log.Fatal(err)
	}

	// Create signal channel for testing
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Create context for testing
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init Core
	pgcnsl := NewPgpoolConsul(c)

	// Run process in background
	go pgcnsl.Run()

	// Run handle shutdown
	go pgcnsl.handleShutdown(signalChan, cancel)

	// Send signal into the channel
	signalChan <- syscall.SIGINT

	// Perform clean up on consul
	cnsl, err := consul.NewConsul(*c)
	if err != nil {
		t.Errorf("Test core: %s", err)
	}

	if err := cnsl.RegisterService("SHUTDOWN"); err != nil {
		t.Errorf("Test core: %s", err)
	}

	// Delay to let the proceess end
	time.Sleep(10 * time.Second)
}
