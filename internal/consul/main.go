package consul

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	capi "github.com/hashicorp/consul/api"
	"github.com/meorkamil/pgpool-consul/internal/model"
)

type consul struct {
	client      *capi.Client
	ServiceName string
	Port        int
	Addr        string
	ServiceId   string
	Interval    string
	Timeout     string
}

func NewConsul(c model.Config) (*consul, error) {
	// Get hostname to be in serviceid
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("error getting hostname: %w", err)
	}

	// Create client connection to Consul
	client, err := capi.NewClient(&capi.Config{
		Address: c.Consul.Addr,
	})
	if err != nil {
		return nil, fmt.Errorf("creating Consul client: %w", err)
	}

	return &consul{
		client:      client,
		ServiceName: c.Consul.Services.Name,
		Addr:        c.Consul.Services.Addr,
		Port:        c.Consul.Services.Port,
		ServiceId:   fmt.Sprintf("%s-%s", c.Consul.Services.Name, strings.ToLower(hostname)),
		Interval:    c.Consul.Services.Interval, Timeout: c.Consul.Services.Timeout,
	}, nil
}

// RegisterService registers the service based on pgpool state (LEADER/STADBY)
func (c *consul) RegisterService(pgStat string) error {
	serviceRegistration := &capi.AgentServiceRegistration{
		ID:      c.ServiceId,
		Name:    c.ServiceName,
		Port:    c.Port,
		Tags:    []string{strings.ToLower(pgStat)},
		Address: c.Addr,
		Check: &capi.AgentServiceCheck{
			TCP:                            fmt.Sprintf("%s:%s", c.Addr, strconv.Itoa(c.Port)),
			Timeout:                        c.Timeout,
			Interval:                       c.Interval,
			DeregisterCriticalServiceAfter: "1m",
		},
	}

	// Selectively register and deregister based on state
	switch {
	case pgStat == "LEADER" || pgStat == "STANDBY":
		return c.registerServiceWithTag(pgStat, serviceRegistration)

	case pgStat == "SHUTDOWN":
		return c.deregisterServiceIfExists()

	default:
		slog.Error(fmt.Sprintf("PGPool Invalid state: %s. Unable to proceed with service registartion %s", pgStat, c.ServiceId))
		return c.deregisterServiceIfExists()
	}
}

// Register service with tag
func (c *consul) registerServiceWithTag(pgStat string, serviceRegistration *capi.AgentServiceRegistration) error {
	err := c.client.Agent().ServiceRegister(serviceRegistration)
	if err != nil {
		return fmt.Errorf("failed to register service %s with tag %s %w", c.ServiceName, pgStat, err)
	}
	slog.Info(fmt.Sprintf("Service %s registered with tag %s", c.ServiceId, pgStat))
	return nil
}

// Deregister service
func (c *consul) deregisterServiceIfExists() error {
	srv, _, err := c.client.Catalog().Service(c.ServiceName, "", nil)
	if err != nil {
		return fmt.Errorf("catalog error: %w", err)
	}

	if len(srv) == 0 {
		slog.Info(fmt.Sprintf("Service %s not found", c.ServiceId))
		return nil
	}

	for _, v := range srv {
		if v.ServiceID == c.ServiceId {
			err := c.client.Agent().ServiceDeregister(c.ServiceId)
			if err != nil {
				return fmt.Errorf("failed to deregister service %s: %w", c.ServiceId, err)
			}
			slog.Info(fmt.Sprintf("Service %s deregistered", c.ServiceId))
		}
	}
	return nil
}
