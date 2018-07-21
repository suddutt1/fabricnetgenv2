package common

import (
	"fmt"
	"io/ioutil"
	"strings"

	util "github.com/suddutt1/fabricnetgenv2/util"

	"gopkg.in/yaml.v2"
)

//GenerateCrytoConfig generate the CryptoConfig file
func GenerateCrytoConfig(nc *NetworkConfig, filePath string) bool {
	useCA := false
	isSuccess := true
	rootConfig := nc.GetRootConfig()
	useCA = util.GetBoolean(rootConfig["addCA"])
	//Perform the orderer part
	ordererConfig := util.GetMap(rootConfig["orderers"])
	if ordererConfig == nil {
		fmt.Println("No orderer specified")
		return false
	}

	cryptoConfig := make(map[string]interface{})
	orderOrgs := make([]map[string]interface{}, 0)
	orderOrgs = append(orderOrgs, buildOrderConfig(ordererConfig))
	cryptoConfig["OrdererOrgs"] = orderOrgs
	orgs, orgsExists := rootConfig["orgs"].([]interface{})
	if !orgsExists {
		fmt.Println("No organizations specified")
		return false
	}
	peerOrgs := make([]map[string]interface{}, 0)
	for _, orgConfig := range orgs {
		peerOrgs = append(peerOrgs, buildOrgConfig(util.GetMap(orgConfig), useCA))
	}
	cryptoConfig["PeerOrgs"] = peerOrgs
	outBytes, _ := yaml.Marshal(cryptoConfig)
	//fmt.Printf("Crypto config Orderes\n%s\n", string(outBytes))

	ioutil.WriteFile(filePath, outBytes, 0666)
	return isSuccess
}
func buildOrderConfig(ordererConfig map[string]interface{}) map[string]interface{} {
	outputStructure := make(map[string]interface{})
	outputStructure["Name"] = util.GetString(ordererConfig["name"])
	outputStructure["Domain"] = util.GetString(ordererConfig["domain"])
	specs := make([]map[string]interface{}, 0)

	//Assuing one as of now
	sansInput := strings.Split(util.GetString(ordererConfig["SANS"]), ",")

	sansArray := make([]string, len(sansInput))
	for indx, sans := range sansInput {
		sansArray[indx] = sans
	}
	sansSpec := make(map[string]interface{})
	sansSpec["SANS"] = sansArray
	specs = append(specs, sansSpec)

	if util.IfEntryExistsInMap(ordererConfig, "haCount") && util.IfEntryExistsInMap(ordererConfig, "type") {
		if util.GetString(ordererConfig["type"]) == "kafka" {
			template := make(map[string]interface{})
			template["Count"] = ordererConfig["haCount"]
			template["Hostname"] = fmt.Sprintf("%s{{.Index}}", util.GetString(ordererConfig["ordererHostname"]))
			outputStructure["Template"] = template
		}
	} else {
		hostnameSpec := make(map[string]interface{})
		hostnameSpec["Hostname"] = util.GetString(ordererConfig["ordererHostname"])
		specs = append(specs, hostnameSpec)
	}
	outputStructure["Specs"] = specs
	return outputStructure
}
func buildOrgConfig(orgConfig map[string]interface{}, useCA bool) map[string]interface{} {
	outputStructure := make(map[string]interface{})
	outputStructure["Name"] = util.GetString(orgConfig["name"])
	outputStructure["Domain"] = util.GetString(orgConfig["domain"])
	if useCA == true {
		caTemplate := make(map[string]string)
		caTemplate["Hostname"] = "ca"
		outputStructure["CA"] = caTemplate
	}

	template := make(map[string]interface{})
	template["Count"] = orgConfig["peerCount"]
	//Assuing one as of now
	sansArray := make([]string, 1)
	sansArray[0] = util.GetString(orgConfig["SANS"])
	template["SANS"] = sansArray
	outputStructure["Template"] = template
	users := make(map[string]interface{})
	users["Count"] = orgConfig["userCount"]
	outputStructure["Users"] = users
	return outputStructure
}
