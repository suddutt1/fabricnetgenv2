package common

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func GenerateConfigForMultipleMachine(nc *NetworkConfig, basePath string) bool {
	var serviceConf ServiceConfig
	serviceConf.Version = "2"
	containers := make(map[string]interface{})
	ordererContainer := BuildOrdererContainer(nc, "./")
	containers[ordererContainer.ContainerName] = containers
	serviceConf.Services = containers
	serviceBytes, _ := yaml.Marshal(serviceConf)
	ioutil.WriteFile(basePath+"/orderer/docker-compose-orderer.yaml", serviceBytes, 0666)
	return true
}
func GenerateComposeYamlFile(containers []Container, filePath string) {

}
