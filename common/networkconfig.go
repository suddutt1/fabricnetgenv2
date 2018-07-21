package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"

	util "github.com/suddutt1/fabricnetgenv2/util"
)

const _VERSION_COMP_MAP = `
{
	"1.0.0":{ "fabricCore":"1.0.0","thirdParty":"1.0.0"},
	"1.0.4":{ "fabricCore":"1.0.4","thirdParty":"1.0.4"},
	"1.1.0":{ "fabricCore":"1.1.0","thirdParty":"1.0.6"}
	
}
`

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
func (nc *NetworkConfig) GetRootConfig() map[string]interface{} {
	return nc.config
}
func (nc *NetworkConfig) DetermineImageVersions() {
	if util.IfEntryExistsInMap(nc.config, "fabricVersion") {
		version, _ := nc.config["fabricVersion"].(string)
		core, thridParty := nc.GetVersions(version)
		nc.config["fabricVersion"] = core
		nc.config["thirdPartyVersion"] = thridParty

	} else {
		nc.config["fabricVersion"] = "1.1.0"
		nc.config["thirdPartyVersion"] = "1.0.6"
	}
}
func (nc *NetworkConfig) GetVersions(version string) (string, string) {
	tmpl, err := template.New("versionMap").Parse(_VERSION_COMP_MAP)
	if err != nil {
		fmt.Printf("Error in reading template %v\n", err)
		return "1.0.0", "1.0.0"
	}
	dataMapContainer := make(map[string]interface{})
	var outputBytes bytes.Buffer
	err = tmpl.Execute(&outputBytes, dataMapContainer)
	if err != nil {
		fmt.Printf("Error in generating the version map file %v\n", err)
		return "1.0.0", "1.0.0"
	}
	versionMap := make(map[string]map[string]string)
	json.Unmarshal(outputBytes.Bytes(), &versionMap)
	if _, isOk := versionMap[version]; !isOk {
		fmt.Println("Invalid version number provided defaulting to 1.0.0")
		return "1.0.0", "1.0.0"
	}
	coreVersion := versionMap[version]["fabricCore"]
	thirdPartyVersion := versionMap[version]["thirdParty"]
	return coreVersion, thirdPartyVersion
}
