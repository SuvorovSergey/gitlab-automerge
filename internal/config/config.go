package config

import (
	"errors"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Gitlab    GitlabConfig
	Scheduler *SchedulerConfig
}

type GitlabConfig struct {
	Token      string `yaml:"token" env:"GITLAB_TOKEN"`
	BaseUrl    string `yaml:"baseUrl" env:"GITLAB_BASE_URL"`
	ConfigFile string `yaml:"configFile" env:"CONFIG_FILE"`
}

type SchedulerConfig struct {
	Merge  time.Duration `yaml:"merge" env:"MERGE"`
	Update time.Duration `yaml:"update" env:"UPDATE"`
}

var cfg Config

func NewConfig(filename string) (*Config, error) {
	err := cleanenv.ReadConfig(filename, &cfg)
	if err != nil {
		return &cfg, err
	}

	err = cfg.validate()

	if err != nil {
		return &cfg, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Gitlab.BaseUrl) == 0 {
		return errors.New("you must provide Gitlab base url")
	}
	if len(c.Gitlab.Token) == 0 {
		return errors.New("you must provide Gitlab Private Token")
	}
	if len(c.Gitlab.ConfigFile) == 0 {
		return errors.New("you must provide Gitlab project config file name")
	}

	return nil
}
