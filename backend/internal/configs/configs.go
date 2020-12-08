package configs

import (
	"encoding/json"
	"os"
)

type ConfI interface {
	GetPort() string
}

type Conf struct {
	Server struct {
		Port string
	}
}

func (c *Conf) GetPort() string {
	return c.Server.Port
}

func Read(confPath string) (*Conf, error) {
	file, err := os.Open(confPath)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	var config Conf
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	return &config, nil
}
