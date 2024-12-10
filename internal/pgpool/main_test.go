package pgpool

import (
	"log"
	"testing"

	"github.com/meorkamil/pgpool-consul/internal/util"
	"github.com/stretchr/testify/assert"
)

const (
	TESTCONFIG = "../../config/config.yml"
)

func TestPgpool(t *testing.T) {
	// Init configuration file
	c, err := util.ConfigInit(TESTCONFIG, "test")
	if err != nil {
		log.Fatal(err)
	}

	// Assert consul
	assert.NotNil(t, c)

	// Assert consul
	assert.Equal(t, "http://localhost:8500", c.Consul.Addr)
	assert.Equal(t, "pgpool", c.Consul.Services.Name)
	assert.Equal(t, "localhost", c.Consul.Services.Addr)
	assert.Equal(t, 9999, c.Consul.Services.Port)
	assert.Equal(t, "5s", c.Consul.Services.Interval)
	assert.Equal(t, "2s", c.Consul.Services.Timeout)

	assert.Equal(t, "localhost", c.Pgpool.Listen)
	assert.Equal(t, "~/home/postgres/.pcppass", c.Pgpool.Pcppassfile)
	assert.Equal(t, 9898, c.Pgpool.Pcpport)
	assert.Equal(t, "pgpool", c.Pgpool.Pcpuser)
	assert.Equal(t, "0", c.Pgpool.Id)

	assert.Equal(t, 3, c.Global.Interval)
	assert.Equal(t, "test", c.Version)

	// Mock channel for pgpool state
	pgpoolChan := make(chan string, 1)
	pgpool := NewPgPool(*c)

	// Start pgpool
	go pgpool.Run(pgpoolChan)

	select {
	// Handle channel once received
	case pgpoolStat := <-pgpoolChan:
		t.Logf("Test pgpool: %s", pgpoolStat)
		close(pgpoolChan)
	}
}
