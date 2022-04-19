package types_test

import (
	"encoding/json"
	"testing"
)



type AppConfig struct {
	Node struct {
		NodeAddress string `json:"node_address"`
	} `json:"node"`

	Server struct {
		GRPCAddress    string `json:"grpc_address"`
		GatewayAddress string `json:"gateway"`
	} `json:"server"`
}

func TestGood(t *testing.T) {
	appConfig := AppConfig{}
	appConfig.Node.NodeAddress = "localhost:9090"
	appConfig.Server.GatewayAddress = "localhost:9094"
	appConfig.Server.GRPCAddress = "localhost:9095"
	data, err := json.Marshal(appConfig)
	if err != nil  {
		t.Log(err)
	}

	t.Log(string(data))


}

