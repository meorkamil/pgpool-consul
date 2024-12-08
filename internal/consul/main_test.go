package consul

import (
	"log"
	"testing"

	"github.com/meorkamil/pgpool-consul/internal/util"
	"github.com/stretchr/testify/assert"
)

const (
	TESTCONFIG = "../../config/config.yml"
)

func TestConsul(t *testing.T) {
	// Init configuration file
	c, err := util.ConfigInit(TESTCONFIG, "test")
	if err != nil {
		log.Fatal(err)
	}

	// Assert consul
	assert.NotNil(t, c)

	// Create consul client
	cnsl, err := NewConsul(*c)
	if err != nil {
		t.Errorf("Test consul %s", err)
	}

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

	// Mock pgpool state
	pgpoolStat := []string{"LEADER", "STANDBY", "ERROR"}

	// Register consul service
	for _, v := range pgpoolStat {
		if err := cnsl.RegisterService(v); err != nil {
			t.Errorf("Test consul registration: %s", err)
		}

	}

	// Deregister consul service
	if err := cnsl.deregisterServiceIfExists(); err != nil {
		t.Errorf("Test consul registration: %s", err)
	}
}
