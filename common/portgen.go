package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	util "github.com/suddutt1/fabricnetgenv2/util"
)

type PortManager struct {
	isContineousPort bool
	portMap          map[string][]string
	allocationMap    map[string]string
	startPort        int
	lastPortAssigned int
}

func (pm *PortManager) Init(nc *NetworkConfig) {
	//Check the networkconfig inputs and initializes various
	//variables
	if util.IfEntryExistsInMap(nc.GetRootConfig(), "startPort") {
		pm.isContineousPort = true
		pm.startPort = util.GetNumber(nc.GetRootConfig()["startPort"])
		pm.lastPortAssigned = pm.startPort
	}
	pm.portMap = make(map[string][]string)
	pm.allocationMap = make(map[string]string)
	pm.generateSequence(7050) //For orderer
	pm.generateSequence(7051) //GRPC URL for peer
	pm.generateSequence(7053) //Event Hub for peer
	pm.generateSequence(5984) //For CouchDB

}
func (pm *PortManager) generateSequence(basePortNumber int) {
	portList := make([]string, 0)
	for port := basePortNumber; port < 32000; port = port + 1000 {
		portList = append(portList, fmt.Sprintf("%d", port))
	}
	pm.portMap[fmt.Sprintf("%d", basePortNumber)] = portList
}
func (pm *PortManager) GetOrdererPort(ordererHostName string) string {

	return pm.allocatePort(ordererHostName, "7050")
}
func (pm *PortManager) GetEventPort(peerHostName string) string {
	return pm.allocatePort(peerHostName, "7053")

}
func (pm *PortManager) GetCouchPort(couchHostName string) string {
	return pm.allocatePort(couchHostName, "5984")
}
func (pm *PortManager) GetGRPCPort(peerHostName string) string {
	return pm.allocatePort(peerHostName, "7051")

}
func (pm *PortManager) allocatePort(hostname, basePort string) string {
	key := fmt.Sprintf("%s:%s", hostname, basePort)
	if pm.isContineousPort {
		returnValue := fmt.Sprintf("%d:%s", pm.lastPortAssigned, basePort)
		pm.allocationMap[key] = fmt.Sprintf("%d", pm.lastPortAssigned)
		pm.lastPortAssigned = pm.lastPortAssigned + 1
		return returnValue
	}
	if availablePortList, isOk := pm.portMap[basePort]; isOk && len(availablePortList) > 0 {
		port := availablePortList[0]
		returnValue := fmt.Sprintf("%s:%s", port, basePort)
		pm.allocationMap[key] = port
		pm.portMap[basePort] = availablePortList[1:]
		return returnValue
	}
	return "######"
}
func (pm *PortManager) PrintAllocationMap(filePath string) bool {
	prettyBytes, _ := json.MarshalIndent(pm.allocationMap, "", " ")
	ioutil.WriteFile(filePath, prettyBytes, 0666)
	return true
}
