package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	util "github.com/suddutt1/fabricnetgenv2/util"
)

func GenerateChainCodeScriptsSingleMachine(nc *NetworkConfig, path string) bool {
	fmt.Println("Generating config scripts")
	fileNames := make([]string, 0)
	dataMapContainer := nc.GetRootConfig()
	//Build the msp info
	mspMap := make(map[string]string)
	peerCountMap := make(map[string]int)
	ordererConfig := util.GetMap(dataMapContainer["orderers"])
	ordererFDQN := util.GetString(ordererConfig["ordererHostname"]) + "." + util.GetString(ordererConfig["domain"])
	if nc.IsKafkaOrderer() {
		ordererFDQN = util.GetString(ordererConfig["ordererHostname"]) + "0." + util.GetString(ordererConfig["domain"])
	}
	orgs, orgsExists := dataMapContainer["orgs"].([]interface{})
	if !orgsExists {
		fmt.Println("No organizations specified")
		return false
	}
	for _, org := range orgs {
		orgConfig := util.GetMap(org)
		name := util.GetString(orgConfig["name"])
		mspId := util.GetString(orgConfig["mspID"])
		mspMap[name] = mspId
		peerCountMap[name] = util.GetNumber(orgConfig["peerCount"])
	}
	chainCodes, chainExists := dataMapContainer["chaincodes"].([]interface{})
	if !chainExists {
		fmt.Println("No chain codes defined")
		return false
	}

	for _, ccInfo := range chainCodes {

		chainCodeConfig := util.GetMap(ccInfo)
		ccID := util.GetString(chainCodeConfig["ccid"])
		version := util.GetString(chainCodeConfig["version"])
		src := util.GetString(chainCodeConfig["src"])
		channelName := fmt.Sprintf("%schannel", strings.ToLower((util.GetString(chainCodeConfig["channelName"]))))
		participants, particpantExists := chainCodeConfig["participants"].([]interface{})
		if !particpantExists {
			fmt.Printf("No participants \n")
			return false
		}
		instShFileName := path + ccID + "_install.sh"
		fileNames = append(fileNames, instShFileName)
		shFileInstall, _ := os.Create(instShFileName)
		shFileInstall.WriteString("#!/bin/bash\n")
		updShFileName := path + ccID + "_update.sh"
		fileNames = append(fileNames, updShFileName)
		shFileUpdateCC, _ := os.Create(updShFileName)
		shFileUpdateCC.WriteString("#!/bin/bash\n")
		shFileUpdateCC.WriteString("if [[ ! -z \"$1\" ]]; then  \n")
		policy := ""
		for _, participant := range participants {
			peerCount := peerCountMap[util.GetString(participant)]
			for index := 0; index < peerCount; index++ {
				lineToWrite := fmt.Sprintf(". setPeer.sh %s peer%d \n", participant, index)
				setChannel := fmt.Sprintf("export CHANNEL_NAME=\"%s\"\n", channelName)
				shFileInstall.WriteString(lineToWrite)
				shFileInstall.WriteString(setChannel)
				shFileUpdateCC.WriteString("\t" + lineToWrite)
				shFileUpdateCC.WriteString(setChannel)
				exeCommand := fmt.Sprintf("peer chaincode install -n %s -v %s -p %s\n", ccID, version, src)
				shFileInstall.WriteString(exeCommand)
				exeUpdCommand := fmt.Sprintf("peer chaincode install -n %s -v %s -p %s\n", ccID, "$1", src)
				shFileUpdateCC.WriteString("\t" + exeUpdCommand)
			}
			policy = policy + ",'" + (mspMap[util.GetString(participant)]) + ".member'"
		}
		runes := []rune(policy)
		finalPolicy := string(runes[1:])
		instCommand := fmt.Sprintf("peer chaincode instantiate -o %s:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C %s -n %s -v %s -c '{\"Args\":[\"init\",\"\"]}' -P \" OR( %s ) \" \n", ordererFDQN, channelName, ccID, version, finalPolicy)
		shFileInstall.WriteString(instCommand)
		updateCommand := fmt.Sprintf("\tpeer chaincode upgrade -o %s:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C %s -n %s -v %s -c '{\"Args\":[\"init\",\"\"]}' -P \" OR( %s ) \" \n", ordererFDQN, channelName, ccID, "$1", finalPolicy)
		shFileUpdateCC.WriteString(updateCommand)
		shFileUpdateCC.WriteString("else\n")
		shFileUpdateCC.WriteString("\techo \". " + ccID + "_updchain.sh  <Version Number>\" \n")
		shFileUpdateCC.WriteString("fi\n")
		shFileInstall.Close()
		shFileUpdateCC.Close()
	}

	//instCommand =
	//     peer chaincode instantiate -o orderer.kg.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $1 -v $2 -c '{"Args":["init",""]}' -P "OR ('RawMaterialDepartmentMSP.member','ManufacturingDepartmentMSP.member','DistributionCenterMSP.member','DistributionCenterMSP.member')"
	for _, fileName := range fileNames {
		os.Chmod(fileName, 0777)
	}
	return true
}

func CreateDefaultCC(nc *NetworkConfig, basePath string) bool {
	os.Mkdir(basePath+"/cli/chaincode", 0777)
	if nc != nil {
		ccDetailsList := nc.GetChaincodeDetails()
		for _, ccDetails := range ccDetailsList {
			ccPath := util.GetString(ccDetails["src"])
			fmt.Printf("\nCreating chaincode path: %s", ccPath)
			os.MkdirAll(basePath+"/cli/chaincode/"+ccPath, 0777)
			ioutil.WriteFile(basePath+"/cli/chaincode/"+ccPath+"/"+"sc_main.go", []byte(_BASE_CHAIN_CODE), 0666)
		}
	} else {
		fmt.Println(_BASE_CHAIN_CODE)
	}
	return true
}

const _BASE_CHAIN_CODE = `
package main

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

/**
  A common base for Hyperledger Fabric Smart Contracts

**/
var _logger = shim.NewLogger("SmartContract.Main")

type SmartContract struct {
}

func (sc *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	_logger.Infof("Inside the main init method")
	return shim.Success([]byte(fmt.Sprintf("{\"isSuccess\": true, \"message\":\"Init successful\"}")))

}
func (sc *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	var response peer.Response
	action, args := stub.GetFunctionAndParameters()
	switch action {
	case "probe":
		response = shim.Success([]byte(fmt.Sprintf("{\"isSuccess\": true, \"ts\":\"%s\"}", time.Now().String())))
	case "save":
		key := args[0]
		value := args[1]
		err := stub.PutState(key, []byte(value))
		if err != nil {
			response = shim.Error(fmt.Sprintf("{\"isSuccess\": false, \"message\":\"Can not save\"}"))
		} else {
			stub.SetEvent("SAVE_EVENT", []byte(value))
			response = shim.Success([]byte(fmt.Sprintf("{\"isSuccess\": true, \"message\":\"Save successful for key %s\"}", key)))
		}
	case "retrieve":
		key := args[0]
		storedValue, err := stub.GetState(key)
		if err != nil {
			response = shim.Error(fmt.Sprintf("{\"isSuccess\": false, \"message\":\"Retrival error\"}"))
		} else {
			response = shim.Success(storedValue)
		}
	default:
		response = shim.Error(fmt.Sprintf("{\"isSuccess\": false, \"message\":\"No action found\"}"))
	}
	return response
}

func main() {
	if err := shim.Start(new(SmartContract)); err != nil {
		_logger.Errorf("Error starting SmartContract chaincode: %s", err)
	}
}


`
