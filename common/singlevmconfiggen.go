package common

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func GenerateConfigForSingleMachine(nc *NetworkConfig, basePath string) bool {
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
	if !GenerateBaseYAML(nc, basePath) {
		fmt.Println("Error in generating base.yaml")
		return false
	}
	if !GenerateSingleMachineDockerFile(nc, basePath) {
		fmt.Println("Could not generate docker-compose.yaml file")
		return false
	}
	if !GenerateGenerateArtifactsScript(nc, basePath+"/generateConfig.sh") {
		fmt.Println("Error in generating generateConfig.sh script ")
		return false
	}
	if !GenerateSetPeerScript(nc, basePath+"/setPeer.sh") {
		return false
	}
	if !GenerateBuildAndJoinChannelScript(nc, basePath+"/setupChannels.sh") {
		return false
	}
	if !GenerateChainCodeScriptsSingleMachine(nc, basePath) {
		return false
	}
	if !GenerateOtherScripts(nc, basePath) {
		return false
	}
	if !GenerateRemoveImagesScript(nc, basePath+"/removeOldImages.sh") {
		return false
	}
	if !GenerateCleanUpScript(nc, basePath+"/cleanup.sh") {
		fmt.Println("Error in generatng cleanup.sh script")
		return false
	}
	nc.PortManager.PrintAllocationMap(basePath + "/portMap.json")
	return true
}

func GenerateSingleMachineDockerFile(nc *NetworkConfig, basePath string) bool {

	addCA := nc.IsCARequired()
	//networkConfig := nc.GetRootConfig()
	var serviceConf ServiceConfig
	serviceConf.Version = "2"
	netWrk := make(map[string]interface{})

	netWrk["fabricnetwork"] = make(map[string]string)
	serviceConf.Network = netWrk
	containers := make(map[string]interface{})
	//Add the orderer
	ordererContainerList := make([]string, 0)
	cliDependencyList := make([]string, 0)
	if !nc.IsKafkaOrderer() {
		orderContainer := BuildOrdererSingleVMSolo(nc, ".")
		containers[orderContainer.ContainerName] = orderContainer
		ordererContainerList = append(ordererContainerList, orderContainer.ContainerName)
		cliDependencyList = append(cliDependencyList, orderContainer.ContainerName)
	}
	//Generate the docker-compose file now
	peerContainers := BuildPeersSingleVM(nc, ordererContainerList)
	//Add the peerContainers into map of containers
	for _, container := range peerContainers {
		containers[container.ContainerName] = container
		cliDependencyList = append(cliDependencyList, container.ContainerName)
	}
	cli := BuildCLIForSingleMachine("./", cliDependencyList)
	containers[cli.ContainerName] = cli
	serviceConf.Services = containers
	serviceBytes, _ := yaml.Marshal(serviceConf)
	if addCA == true {
		ioutil.WriteFile(basePath+"/docker-compose-template.yaml", serviceBytes, 0666)
	} else {
		ioutil.WriteFile(basePath+"/docker-compose.yaml", serviceBytes, 0666)
	}

	return true

}
