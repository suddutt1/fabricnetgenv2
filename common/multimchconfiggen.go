package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	util "github.com/suddutt1/fabricnetgenv2/util"
	yaml "gopkg.in/yaml.v2"
)

func GenerateConfigForMultipleMachine(nc *NetworkConfig, basePath string) bool {

	nc.Init()
	commonArtifactsPath := basePath + "/common"
	if !GenerateDownloadScripts(nc, commonArtifactsPath) {
		fmt.Println("Error in generating download scritps")
		return false
	}
	if !GenerateConfigTxGen(nc, commonArtifactsPath+"/configtx.yaml") {
		fmt.Println("Error in generating configtx.yaml file")
		return false
	}

	if !GenerateCrytoConfig(nc, commonArtifactsPath+"/crypto-config.yaml") {
		fmt.Println("Error in generating crypto-config.yaml")
		return false
	}
	if !GenerateBaseYAML(nc, commonArtifactsPath) {
		fmt.Println("Error in generating base.yaml")
		return false
	}
	if !GenerateMultiMachineOrderer(nc, basePath) {
		return false
	}
	if !GenerateMultiMachinePeers(nc, basePath) {
		return false
	}
	if !GenerateGenerateArtifactsScript(nc, commonArtifactsPath+"/generateConfig.sh") {
		fmt.Println("Error in generating generateConfig.sh script ")
		return false
	}

	if !GenerateSetPeerScript(nc, commonArtifactsPath+"/setPeer.sh") {
		return false
	}
	if !GenerateBuildAndJoinChannelScript(nc, commonArtifactsPath+"/setupChannels.sh") {
		return false
	}
	if !GenerateChainCodeScriptsSingleMachine(nc, commonArtifactsPath) {
		return false
	}

	if !GenerateOtherScripts(nc, commonArtifactsPath) {
		return false
	}
	if !GenerateRemoveImagesScript(nc, commonArtifactsPath+"/removeOldImages.sh") {
		return false
	}
	if !GenerateCleanUpScript(nc, commonArtifactsPath+"/cleanup.sh") {
		fmt.Println("Error in generatng cleanup.sh script")
		return false
	}
	GenerateDistributeConfig(nc, basePath)

	return true
}
func GenerateDistributeConfig(nc *NetworkConfig, basePath string) {

}
func GenerateMultiMachineOrderer(nc *NetworkConfig, basePath string) bool {
	ordererContainer := BuildOrdererContainer(nc, ".")
	ordererBaseDir := basePath + "/orderer/"

	if !createDir(ordererBaseDir) {

		return false
	}
	GenerateComposeYamlFile([]Container{ordererContainer}, ordererBaseDir+"docker-compose.yaml")
	return true
}

func GenerateMultiMachinePeers(nc *NetworkConfig, basePath string) bool {
	orgConfigs, _ := nc.GetRootConfig()["orgs"].([]interface{})
	//containers := make([]Container, 0)
	couchCount := 0
	for _, org := range orgConfigs {
		orgConfig, _ := org.(map[string]interface{})
		fmt.Printf("Processing org %s \n", orgConfig["name"])
		peerCountFlt, _ := orgConfig["peerCount"].(float64)
		peerCount := int(peerCountFlt)
		fmt.Printf("\tPeer count is %d \n ", peerCount)
		//TODO: AddCA
		for peerIndex := 0; peerIndex < peerCount; peerIndex++ {
			peerID := fmt.Sprintf("peer%d", peerIndex)
			couchID := fmt.Sprintf("couch%d", couchCount)
			couchContainer := BuildCouchDB(couchID, nc)
			peerContainer := BuildPeerImage(".", peerID, util.GetString(orgConfig["domain"]), util.GetString(orgConfig["mspID"]), couchID, []string{}, nc)
			//containers = append(containers, couchContainer)
			//containers = append(containers, peerContainer)
			orgName, _ := orgConfig["name"].(string)
			peerDir := fmt.Sprintf("%s/%s-%s/", basePath, peerID, strings.ToLower(orgName))
			couchDir := fmt.Sprintf("%s/%s-couch-%s/", basePath, peerID, strings.ToLower(orgName))
			if !createDir(peerDir) || !createDir(couchDir) {
				return false
			}

			GenerateComposeYamlFile([]Container{peerContainer}, peerDir+"docker-compose.yaml")
			GenerateComposeYamlFile([]Container{couchContainer}, couchDir+"docker-compose.yaml")
			couchCount++

		}
	}
	cliDir := basePath + "/cli/"

	cliContainer := BuildCLIContainer("./", []string{}, nc)
	if !createDir(cliDir) {
		return false
	}
	GenerateComposeYamlFile([]Container{cliContainer}, cliDir+"docker-compose.yaml")
	return true
}
func createDir(dirPath string) bool {
	err := os.MkdirAll(dirPath, 0777)
	if err != nil {
		fmt.Printf("\nUnable to generate directory %+v\n", err)
		return false
	}
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
