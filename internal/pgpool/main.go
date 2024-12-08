package pgpool

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/meorkamil/pgpool-consul/internal/model"
)

// Constants for different states
const (
	WATCHDOGCMD  = "pcp_watchdog_info"
	LeaderState  = "LEADER"
	StandbyState = "STANDBY"
	UnknownState = "UNKNOWN"
	ErrorState   = "ERROR"
)

type pgpool struct {
	Pcppassfile string
	Timeout     time.Duration
	Pcpuser     string
	Pcpport     string
	Pcpaddr     string
	Id          string
}

func NewPgPool(c model.Config) *pgpool {
	return &pgpool{
		Pcppassfile: c.Pgpool.Pcppassfile,
		Timeout:     3 * time.Second,
		Pcpuser:     c.Pgpool.Pcpuser,
		Pcpport:     strconv.Itoa(c.Pgpool.Pcpport),
		Pcpaddr:     c.Pgpool.Listen,
		Id:          c.Pgpool.Id,
	}
}

// Run starts the pgpool service and sends its state to the provided channel
func (c *pgpool) Run(ch chan string) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	// Get the output from the command
	cmdOut, err := c.cmd(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
		ch <- ErrorState
		return
	}

	// Send the state to the channel based on the command output
	switch cmdOut {
	case LeaderState:
		ch <- LeaderState

	case StandbyState:
		ch <- StandbyState

	default:
		ch <- UnknownState
	}
}

// cmd executes the given command with the provided context and returns the output
func (c *pgpool) cmd(ctx context.Context) (string, error) {
	// Set command for pcp_watchdog_info
	PCPCMD := fmt.Sprintf(
		"PCPPASSFILE=%s %s -h %s -U%s -p %s -n %s -v | grep \"Status Name\" | awk '{print $4}'",
		c.Pcppassfile, WATCHDOGCMD, c.Pcpaddr, c.Pcpuser, c.Pcpport, c.Id,
	)

	// Execute command
	out, err := exec.CommandContext(ctx, "bash", "-c", PCPCMD).Output()

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("PCP command %w", ctx.Err())
		}
		return "", fmt.Errorf("PCP command failed: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}
