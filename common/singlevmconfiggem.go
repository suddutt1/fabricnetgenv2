package common

import "fmt"

func GenerateConfigForSingleVM(nc *NetworkConfig, basePath string) bool {

	if !GenerateDownloadScripts(nc, basePath) {
		fmt.Println("Error in generating download scritps")
		return false
	}
	if !GenerateCrytoConfig(nc, basePath+"/crypto-config.yaml") {
		fmt.Println("Error in generating crypto-config.yaml")
		return false
	}
	return true
}
