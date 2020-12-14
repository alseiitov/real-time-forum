package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ConfI interface {
	GetPort() string
	GetAPIAdress() string
}

type Conf struct {
	Server struct {
		Port string `json:"port"`
	} `json:""server"`
	API struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"api"`
	InJSON string
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
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var config Conf
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	config.InJSON = string(data)
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
