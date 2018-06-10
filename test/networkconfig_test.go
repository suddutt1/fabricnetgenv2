package test

import (
	"encoding/json"
	"testing"

	common "github.com/suddutt1/fabricnetgenv2/common"
	util "github.com/suddutt1/fabricnetgenv2/util"
)

func Test_NetworkUnmarshal(t *testing.T) {
	t.Logf("Testing networkconfig")
	_, networkConfigBytes := util.LoadTemplate(_network_config_1, "network-config", nil)
	netConfig := common.NetworkConfig{}
	err := json.Unmarshal(networkConfigBytes, &netConfig)
	if err != nil {
		t.FailNow()
	}
	netConfig.PrintConfig()
}

const _network_config_1 = `
{
	"a":"b"
}
`
