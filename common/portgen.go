package common

import "fmt"

type PortManager struct {
	isContineousPort bool
	portMap          map[string][]string
	allocationMap    map[string]string
	startPort        int
	lastPort         int
}

func (pm *PortManager) Init(nc *NetworkConfig) {
	//Check the networkconfig inputs and initializes various
	//variables
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
	if availablePortList, isOk := pm.portMap[basePort]; isOk && len(availablePortList) > 0 {
		port := availablePortList[0]
		returnValue := fmt.Sprintf("%s:%s", port, basePort)
		pm.allocationMap[key] = port
		pm.portMap[basePort] = availablePortList[1:]
		return returnValue
	}
	return "######"
}
