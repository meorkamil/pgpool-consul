package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	// Create a struct value
	c := Config{
		Consul: struct {
			Addr     string "yaml:\"addr\""
			Services struct {
				Name     string "yaml:\"name\""
				Addr     string "yaml:\"addr\""
				Port     int    "yaml:\"port\""
				Interval string "yaml:\"interval\""
				Timeout  string "yaml:\"timeout\""
			} "yaml:\"services\""
		}{
			Addr: "http://localhost:8500",
			Services: struct {
				Name     string "yaml:\"name\""
				Addr     string "yaml:\"addr\""
				Port     int    "yaml:\"port\""
				Interval string "yaml:\"interval\""
				Timeout  string "yaml:\"timeout\""
			}{
				Name:     "pgpool",
				Addr:     "localhost",
				Port:     9999,
				Interval: "5s",
				Timeout:  "2s",
			},
		},
		Pgpool: struct {
			Listen      string "yaml:\"listen\""
			Pcppassfile string "yaml:\"pcppassfile\""
			Pcpport     int    "yaml:\"pcpport\""
			Pcpuser     string "yaml:\"pcpuser\""
			Id          string "yaml:\"id\""
		}{
			Listen:      "localhost",
			Pcppassfile: "~/home/postgres/.pcppass",
			Pcpport:     9898,
			Pcpuser:     "pgpool",
			Id:          "0",
		},
		Global: struct {
			Interval int "yaml:\"interval\""
		}{
			Interval: 3,
		},
		Version: "test",
	}

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
}
