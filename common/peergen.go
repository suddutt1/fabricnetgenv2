package common

import (
	"fmt"

	util "github.com/suddutt1/fabricnetgenv2/util"
)

func BuildPeersSingleVM(nc *NetworkConfig, otherDependenies []string) []Container {
	orgConfigs, _ := nc.GetRootConfig()["orgs"].([]interface{})
	containers := make([]Container, 0)
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
			peerContainer := BuildPeerImage(".", peerID, util.GetString(orgConfig["domain"]), util.GetString(orgConfig["mspID"]), couchID, otherDependenies, nc)
			containers = append(containers, couchContainer)
			containers = append(containers, peerContainer)
			couchCount++

		}
	}
	return containers
}

func BuildPeerImage(cryptoBasePath, peerId, domainName, mspID, couchID string, otherDependencies []string, nc *NetworkConfig) Container {

	extnds := make(map[string]string)
	extnds["file"] = "base.yaml"
	extnds["service"] = "peer"
	peerFQDN := peerId + "." + domainName

	peerEnvironment := make([]string, 0)
	peerEnvironment = append(peerEnvironment, "CORE_PEER_ID="+peerFQDN)
	peerEnvironment = append(peerEnvironment, "CORE_PEER_ADDRESS="+peerFQDN+":7051")
	peerEnvironment = append(peerEnvironment, "CORE_PEER_CHAINCODELISTENADDRESS="+peerFQDN+":7052")
	peerEnvironment = append(peerEnvironment, "CORE_PEER_GOSSIP_EXTERNALENDPOINT="+peerFQDN+":7051")
	peerEnvironment = append(peerEnvironment, "CORE_PEER_LOCALMSPID="+mspID)
	peerEnvironment = append(peerEnvironment, "CORE_LEDGER_STATE_STATEDATABASE=CouchDB")
	peerEnvironment = append(peerEnvironment, "CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS="+couchID+":5984")
	if peerId == "peer0" {
		peerEnvironment = append(peerEnvironment, "CORE_PEER_GOSSIP_BOOTSTRAP="+peerFQDN+":7051")
	} else {
		peerEnvironment = append(peerEnvironment, "CORE_PEER_GOSSIP_BOOTSTRAP=peer0."+domainName+":7051")
	}
	vols := make([]string, 0)
	vols = append(vols, "/var/run/:/host/var/run/")
	vols = append(vols, cryptoBasePath+"/crypto-config/peerOrganizations/"+domainName+"/peers/"+peerFQDN+"/msp:/etc/hyperledger/fabric/msp")
	vols = append(vols, cryptoBasePath+"/crypto-config/peerOrganizations/"+domainName+"/peers/"+peerFQDN+"/tls:/etc/hyperledger/fabric/tls")
	var depends = make([]string, 0)
	depends = append(depends, couchID)
	depends = append(depends, otherDependencies...)
	var networks = make([]string, 0)
	networks = append(networks, "fabricnetwork")

	var container Container
	container.ContainerName = peerFQDN
	container.Environment = peerEnvironment
	container.Volumns = vols
	container.Depends = depends
	container.Networks = networks
	ports := make([]string, 0)
	ports = append(ports, nc.PortManager.GetGRPCPort(peerFQDN))
	ports = append(ports, nc.PortManager.GetEventPort(peerFQDN))

	container.Ports = ports
	container.Extends = extnds

	return container
}

func BuildCouchDB(couchID string, nc *NetworkConfig) Container {
	var couchContainer Container
	couchContainer.ContainerName = couchID
	extnds := make(map[string]string)
	extnds["file"] = "base.yaml"
	extnds["service"] = "couchdb"
	couchContainer.Extends = extnds
	var networks = make([]string, 0)
	networks = append(networks, "fabricnetwork")

	couchContainer.Networks = networks
	ports := make([]string, 0)
	ports = append(ports, nc.PortManager.GetCouchPort(couchID))
	couchContainer.Ports = ports
	return couchContainer
}
