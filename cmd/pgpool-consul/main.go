package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"

	"github.com/meorkamil/pgpool-consul/internal/core"
	"github.com/meorkamil/pgpool-consul/internal/util"
)

const version = "v0.0.0"

func main() {
	// Get config from cmd args
	configPath := flag.String("config", "../../config/config.yml", "full path to configuration file")
	flag.Parse()

	// Parse config
	c, err := util.ConfigInit(*configPath, version)
	if err != nil {
		log.Fatal(err)
	}

	// Start pgpool-consul
	pgcnsl := core.NewPgpoolConsul(c)
	if err := pgcnsl.Run(); err != nil {
		slog.Error(fmt.Sprintf("Failed to start: %s", err))
	}
}
