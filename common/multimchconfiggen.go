package common

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func GenerateConfigForMultipleMachine(nc *NetworkConfig, basePath string) bool {

	nc.Init()
	if !GenerateDownloadScripts(nc, basePath) {
		fmt.Println("Error in generating download scritps")
		return false
	}
	if !GenerateConfigTxGen(nc, basePath+"/configtx.yaml") {
		fmt.Println("Error in generating configtx.yaml file")
		return false
	}

	if !GenerateCrytoConfig(nc, basePath+"/crypto-config.yaml") {
		fmt.Println("Error in generating crypto-config.yaml")
		return false
	}
	if !GenerateGenerateArtifactsScript(nc, basePath+"/generateConfig.sh") {
		fmt.Println("Error in generating generateConfig.sh script ")
		return false
	}
	if !GenerateMultiMachineOrderer(nc, basePath) {
		return false
	}
	if !GenerateOtherScripts(nc, basePath) {
		return false
	}
	return true
}
func GenerateMultiMachineOrderer(nc *NetworkConfig, basePath string) bool {
	ordererContainer := BuildOrdererContainer(nc, ".")
	ordererBaseDir := basePath + "/orderer/"
	err := os.MkdirAll(ordererBaseDir, 0777)
	if err != nil {
		fmt.Printf("\nUnable to generate orderer directory %+v\n", err)
		return false
	}
	GenerateBaseYAML(nc, ordererBaseDir)
	GenerateComposeYamlFile([]Container{ordererContainer}, ordererBaseDir+"docker-compose.yaml")
	return true
}
func GenerateComposeYamlFile(containers []Container, filePath string) {
	var serviceConf ServiceConfig
	serviceConf.Version = "2"
	containerMap := make(map[string]interface{})
	for _, container := range containers {
		containerMap[container.ContainerName] = container
	}

	serviceConf.Services = containerMap
	serviceBytes, _ := yaml.Marshal(serviceConf)
	ioutil.WriteFile(filePath, serviceBytes, 0666)

}
