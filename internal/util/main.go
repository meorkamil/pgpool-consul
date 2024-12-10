package util

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/meorkamil/pgpool-consul/internal/model"
	"github.com/spf13/viper"
)

// Set configuration
func ConfigInit(c string, v string) (*model.Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(c)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Util %s", err)
	}

	var config model.Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("Util %s", err)
	}

	config.Version = v

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	return &config, nil
}
