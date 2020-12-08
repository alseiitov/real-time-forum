package configs

import (
	"encoding/json"
	"os"
	"fmt"
)

type ConfI interface {
	GetPort() string
	GetAPIAdress() string
}

type Conf struct {
	Server struct {
		Port string
	}
	API struct {
		Host string
		Port string
	}
}

func (c *Conf) GetPort() string {
	return c.Server.Port
}

func (c *Conf) GetAPIAdress() string {
	return fmt.Sprintf("%v:%v", c.API.Host, c.API.Port)
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
		config.Server.Port = "8081"
	}
	if config.API.Host == "" {
		config.Server.Port = "localhost"
	}
	if config.API.Port == "" {
		config.Server.Port = "8080"
	}
	return &config, nil
}
