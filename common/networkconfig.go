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

type ServiceConfig struct {
	Version  string                 `yaml:"version,flow"`
	Network  map[string]interface{} `yaml:"networks,omitempty"`
	Services map[string]interface{} `yaml:"services"`
}

type Container struct {
	Image         string            `yaml:"image,omitempty"`
	Restart       string            `yaml:"restart,omitempty"`
	ContainerName string            `yaml:"container_name,omitempty"`
	TTY           bool              `yaml:"tty,omitempty"`
	Extends       map[string]string `yaml:"extends,omitempty"`
	Environment   []string          `yaml:"environment,omitempty"`
	WorkingDir    string            `yaml:"working_dir,omitempty"`
	Command       string            `yaml:"command,omitempty"`
	Volumns       []string          `yaml:"volumes,omitempty"`
	Ports         []string          `yaml:"ports,omitempty"`
	Depends       []string          `yaml:"depends_on,omitempty"`
	Networks      []string          `yaml:"networks,omitempty"`
	NetworkMode   string            `yaml:"network_mode,omitempty"`
}

//NetworkConfig
type NetworkConfig struct {
	config      map[string]interface{}
	PortManager *PortManager
}

func (nc *NetworkConfig) Init() {
	nc.PortManager = new(PortManager)
	nc.PortManager.Init(nc)
	nc.DetermineImageVersions()
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
func (nc *NetworkConfig) GetOrdererMSP() string {
	ordConfig := nc.GetOrderConfig()
	return util.GetString(ordConfig["mspID"])
}
func (nc *NetworkConfig) GetOrderConfig() map[string]interface{} {
	ordConfigInput, _ := nc.config["orderers"].(interface{})
	ordConfig, _ := ordConfigInput.(map[string]interface{})
	return ordConfig
}
func (nc *NetworkConfig) IsCARequired() bool {
	return util.GetBoolean(nc.config["addCA"])
}
func (nc *NetworkConfig) IsKafkaOrderer() bool {
	return false
}
func (nc *NetworkConfig) GetChaincodeDetails() []map[string]interface{} {
	ccDetails := nc.GetRootConfig()["chaincodes"]
	ccArray, _ := ccDetails.([]interface{})
	ccDetailsMapArray := make([]map[string]interface{}, 0)
	for _, ccDetails := range ccArray {
		//fmt.Printf("\n%+v", ccDetails)
		ccDetailsMap, _ := ccDetails.(map[string]interface{})
		ccDetailsMapArray = append(ccDetailsMapArray, ccDetailsMap)

	}
	return ccDetailsMapArray
}
func (nc *NetworkConfig) IsMultiMachine() bool {
	return util.GetBoolean(nc.GetRootConfig()["multiMachine"])
}
