package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
	_ "gopkg.in/yaml.v2"
)

type Config struct {
	TmuxPrefix        string `yaml:"tmux-prefix"`
	RAMMin            int    `yaml:"ram-min"`
	RAMMax            int    `yaml:"ram-max"`
	RootDir           string `yaml:"root-dir"`
	JarName           string `yaml:"jar-name"`
	StartTmuxAttached bool   `yaml:"start-tmux-attached"`
	JavaParams        string `yaml:"java-params"`
}

func (c *Config) loadConfigFromFile(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err := yaml.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}
}

func NewConfig() *Config {
	return new(Config)
}
