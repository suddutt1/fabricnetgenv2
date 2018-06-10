package common

import (
	"encoding/json"
	"fmt"
)

//NetworkConfig
type NetworkConfig struct {
	config map[string]interface{}
}

func (nc *NetworkConfig) UnmarshalJSON(data []byte) error {
	nc.config = make(map[string]interface{})
	return json.Unmarshal(data, &nc.config)
}
func (nc NetworkConfig) PrintConfig() {
	fmt.Printf("\n%+v\n", nc.config)
}
