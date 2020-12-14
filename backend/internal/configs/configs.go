package configs

import (
	"encoding/json"
	"fmt"
	"os"
)

type Conf struct {
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
	Database struct {
		Driver   string `json:"driver"`
		Path     string `json:"path"`
		FileName string `json:"fileName"`
		Schema   string `json:"schema"`
	} `json:"database"`
}

func (c *Conf) GetPort() string {
	return c.Server.Port
}

func (c *Conf) GetDBDirPath() string {
	return c.Database.Path
}

func (c *Conf) GetDBFilePath() string {
	return fmt.Sprintf("%v/%v", c.Database.Path, c.Database.FileName)
}

func (c *Conf) GetDBDriver() string {
	return c.Database.Driver
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
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	return &config, nil
}
