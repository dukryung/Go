package types

import (
	"encoding/json"
	"io/ioutil"
)

const DefaultAppConfigPath = "./default_app_config_path"

type AppConfig struct {
	Node struct {
		NodeAddress string `json:"node_address"`
	} `json:"node"`

	Server struct {
		GRPCAddress    string `json:"grpc_address"`
		GatewayAddress string `json:"gateway"`
	} `json:"server"`
}

func (config *AppConfig) LoadAppConfig(configPath string) error {
	if configPath == "" {
		configPath = DefaultAppConfigPath
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = config.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	return nil
}

func (config *AppConfig) UnmarshalJSON(data []byte) error {
	type tempAppConfig AppConfig

	appConfig := tempAppConfig{}

	err := json.Unmarshal(data, &appConfig)
	if err != nil {
		return err
	}

	*config = AppConfig(appConfig)

	return nil
}
