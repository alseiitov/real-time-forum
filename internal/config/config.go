package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type Conf struct {
	API      API      `json:"api"`
	Client   Client   `json:"client"`
	Database Database `json:"database"`
	Auth     Auth     `json:"auth"`
	Forum    Forum    `json:"forum"`
}

type API struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Auth struct {
	AccessTokenTTL  int `json:"accessTokenTTL"`
	RefreshTokenTTL int `json:"refreshTokenTTL"`
}

type Client struct {
	Port string `json:"port"`
}

type Database struct {
	Driver     string `json:"driver"`
	Path       string `json:"path"`
	FileName   string `json:"fileName"`
	ImagesDir  string `json:"imagesDir"`
	SchemesDir string `json:"schemesDir"`
}

type Forum struct {
	DefaultMaleAvatar        string `json:"defaultMaleAvatar"`
	DefaultFemaleAvatar      string `json:"defaultFemaleAvatar"`
	PostsForPage             int    `json:"postsForPage"`
	PostsModerationIsEnabled bool   `json:"postsModerationIsEnabled"`
}

func NewConfig(confPath string) (*Conf, error) {
	var config Conf

	file, err := os.Open(confPath)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Conf) GetBackendPort() string {
	return c.API.Port
}

func (c *Conf) GetBackendAdress() string {
	host := c.API.Host
	port := c.GetBackendPort()
	if host == "localhost" || host == "127.0.0.1" {
		return fmt.Sprintf("%s:%s", host, port)
	}
	return host
}

func (c *Conf) GetFrontendPort() string {
	return c.Client.Port
}

func (c *Conf) GetDBFileName() string {
	return c.Database.FileName
}

func (c *Conf) GetDBPath() string {
	return c.Database.Path
}

func (c *Conf) GetDBSchemesDir() string {
	return c.Database.SchemesDir
}

func (c *Conf) GetImagesDir() string {
	return c.Database.ImagesDir
}

func (c *Conf) GetDBDriver() string {
	return c.Database.Driver
}

func (c *Conf) GetTokenTTLs() (time.Duration, time.Duration, error) {
	accessTokenTTL := minutesToDuration(c.Auth.AccessTokenTTL)
	if accessTokenTTL == 0 {
		return 0, 0, errors.New("accessTokenTTL cannot be empty")
	}

	refreshTokenTTL := minutesToDuration(c.Auth.RefreshTokenTTL)
	if refreshTokenTTL == 0 {
		return 0, 0, errors.New("refreshTokenTTL cannot be empty")
	}

	return accessTokenTTL, refreshTokenTTL, nil
}

func (c *Conf) GetDefaultAvatars() (string, string) {
	return c.Forum.DefaultMaleAvatar, c.Forum.DefaultFemaleAvatar
}

func (c *Conf) GetPostsForPage() int {
	return c.Forum.PostsForPage
}

func minutesToDuration(m int) time.Duration {
	return time.Duration(int(time.Minute) * m)
}
