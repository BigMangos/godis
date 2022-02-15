package tcp

import (
	"fmt"
	"time"

	"github.com/hdt3213/godis/config"
)

type Config struct {
	Address    string        `yaml:"address"`
	MaxConnect uint32        `yaml:"max-connect"`
	Timeout    time.Duration `yaml:"timeout"`
}

func InitServerConfig() *Config {
	if config.Properties == nil {
		return nil
	}

	return &Config{
		Address:    fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
		MaxConnect: 0,
		Timeout:    0,
	}
}
