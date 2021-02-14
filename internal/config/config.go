package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Conf struct {
	Backend struct {
		Server struct {
			Host string `json:"host"`
			Port string `json:"port"`
		} `json:"server"`
		Database struct {
			Driver     string `json:"driver"`
			Path       string `json:"path"`
			FileName   string `json:"fileName"`
			SchemesDir string `json:"schemesDir"`
		} `json:"database"`
		Auth struct {
			AccessTokenTTL  int `json:"accessTokenTTL"`
			RefreshTokenTTL int `json:"refreshTokenTTL"`
		} `json:"auth"`
	} `json:"backend"`
	Frontend struct {
		Server struct {
			Port string `json:"port"`
		} `json:"server"`
	} `json:"frontend"`
}

func (c *Conf) GetBackendPort() string {
	return c.Backend.Server.Port
}

func (c *Conf) GetBackendAdress() string {
	host := c.Backend.Server.Host
	port := c.GetBackendPort()
	if host == "localhost" || host == "127.0.0.1" {
		return fmt.Sprintf("%s:%s", host, port)
	}
	return host
}

func (c *Conf) GetFrontendPort() string {
	return c.Frontend.Server.Port
}

func (c *Conf) GetDBFilePath() string {
	return filepath.Join(c.Backend.Database.Path, c.Backend.Database.FileName)
}

func (c *Conf) GetDBSchemesDir() string {
	return c.Backend.Database.SchemesDir
}

func (c *Conf) GetDBDriver() string {
	return c.Backend.Database.Driver
}

func (c *Conf) GetTokenTTLs() (time.Duration, time.Duration, error) {
	accessTokenTTL := minutesToDuration(c.Backend.Auth.AccessTokenTTL)
	if accessTokenTTL == 0 {
		return 0, 0, errors.New("accessTokenTTL cannot be empty")
	}

	refreshTokenTTL := minutesToDuration(c.Backend.Auth.RefreshTokenTTL)
	if refreshTokenTTL == 0 {
		return 0, 0, errors.New("refreshTokenTTL cannot be empty")
	}

	return accessTokenTTL, refreshTokenTTL, nil
}

func minutesToDuration(m int) time.Duration {
	return time.Duration(int(time.Minute) * m)
}

func NewConfig(confPath string) (*Conf, error) {
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
	if config.Backend.Server.Port == "" {
		config.Backend.Server.Port = "8081"
	}
	if config.Frontend.Server.Port == "" {
		config.Frontend.Server.Port = "8080"
	}
	return &config, nil
}
