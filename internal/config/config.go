package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/alseiitov/validator"
)

type Conf struct {
	API      API      `json:"api"`
	Client   Client   `json:"client"`
	Database Database `json:"database"`
	Auth     Auth     `json:"auth"`
	Forum    Forum    `json:"forum"`
}

type API struct {
	Host string `json:"host"	validator:"required"`
	Port string `json:"port"	validator:"required"`
}

type Auth struct {
	AccessTokenTTL  int `json:"accessTokenTTL"	validator:"required"`
	RefreshTokenTTL int `json:"refreshTokenTTL"	validator:"required"`
}

type Client struct {
	Port string `json:"port" validator:"required"`
}

type Database struct {
	Driver     string `json:"driver"			validator:"required"`
	Path       string `json:"path"				validator:"required"`
	FileName   string `json:"fileName"			validator:"required"`
	ImagesDir  string `json:"imagesDir"			validator:"required"`
	SchemesDir string `json:"schemesDir"		validator:"required"`
}

type Forum struct {
	DefaultMaleAvatar              string `json:"defaultMaleAvatar"		validator:"required"`
	DefaultFemaleAvatar            string `json:"defaultFemaleAvatar"	validator:"required"`
	PostsForPage                   int    `json:"postsForPage"			validator:"required"`
	CommentsForPage                int    `json:"commentsForPage"		validator:"required"`
	PostsPreModerationIsEnabled    bool   `json:"postsPreModerationIsEnabled"`
	CommentsPreModerationIsEnabled bool   `json:"commentsPreModerationIsEnabled"`
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

	err = validator.Validate(config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Conf) BackendPort() string {
	return c.API.Port
}

func (c *Conf) BackendAdress() string {
	host := c.API.Host
	port := c.BackendPort()
	if host == "localhost" || host == "127.0.0.1" {
		return fmt.Sprintf("%s:%s", host, port)
	}
	return host
}

func (c *Conf) FrontendPort() string {
	return c.Client.Port
}

func (c *Conf) DBFileName() string {
	return c.Database.FileName
}

func (c *Conf) DBPath() string {
	return c.Database.Path
}

func (c *Conf) DBSchemesDir() string {
	return c.Database.SchemesDir
}

func (c *Conf) ImagesDir() string {
	return c.Database.ImagesDir
}

func (c *Conf) DBDriver() string {
	return c.Database.Driver
}

func (c *Conf) TokenTTLs() (time.Duration, time.Duration, error) {
	accessTokenTTL := minutesToDuration(c.Auth.AccessTokenTTL)
	refreshTokenTTL := minutesToDuration(c.Auth.RefreshTokenTTL)

	return accessTokenTTL, refreshTokenTTL, nil
}

func (c *Conf) DefaultAvatars() (string, string) {
	return c.Forum.DefaultMaleAvatar, c.Forum.DefaultFemaleAvatar
}

func (c *Conf) PostsForPage() int {
	return c.Forum.PostsForPage
}

func (c *Conf) CommentsForPage() int {
	return c.Forum.CommentsForPage
}

func (c *Conf) PostsPreModerationIsEnabled() bool {
	return c.Forum.PostsPreModerationIsEnabled
}

func (c *Conf) CommentsPreModerationIsEnabled() bool {
	return c.Forum.CommentsPreModerationIsEnabled
}

func minutesToDuration(m int) time.Duration {
	return time.Duration(int(time.Minute) * m)
}
