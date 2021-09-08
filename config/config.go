package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Twitter struct {
		AccessToken       string `yaml:"access_token"`
		AccessTokenSecret string `yaml:"access_secret"`
		APIKey            string `yaml:"api_key"`
		APISecretKey      string `yaml:"api_secret"`
	} `yaml:"twitter"`

	Lists struct {
		PositiveIDs []int64 `yaml:"positive"`
		NegativeIDs []int64 `yaml:"negative"`
	} `yaml:"lists"`
}

func Parse(filename string) (c Config, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&c)
	if err != nil {
		return
	}

	return
}
