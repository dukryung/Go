package config

import (
	"encoding/json"
	"io/ioutil"
)

const DefaultConfig = "./config.json"

type AppConfig struct {
	Test ServerConfig
}

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func (config *AppConfig) LoadConfig(path string) error  {
	if path == "" {
		path = DefaultConfig
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		data, err  = ioutil.ReadFile(DefaultConfig)
		if err != nil {
			panic(err)
		}
	}

	err  = json.Unmarshal(data, &AppConfig{})
	if err != nil {
		panic(err)
	}

	return nil

}



