package common

import util "github.com/suddutt1/fabricnetgenv2/util"

//BuildOrdererContainer build the solo order container
func BuildOrdererContainer(nc *NetworkConfig, cryptoBasePath string) Container {
	ordererName := util.GetString(nc.GetOrderConfig()["ordererHostname"])
	domainName := util.GetString(nc.GetOrderConfig()["domain"])
	extnds := make(map[string]string)
	extnds["file"] = "base.yaml"
	extnds["service"] = "orderer"
	ordFQDN := ordererName + "." + domainName
	vols := make([]string, 0)
	vols = append(vols, cryptoBasePath+"/genesis.block:/var/hyperledger/orderer/genesis.block")
	vols = append(vols, cryptoBasePath+"/crypto-config/ordererOrganizations/"+domainName+"/orderers/"+ordFQDN+"/msp:/var/hyperledger/orderer/msp")
	vols = append(vols, cryptoBasePath+"/crypto-config/ordererOrganizations/"+domainName+"/orderers/"+ordFQDN+"/tls/:/var/hyperledger/orderer/tls")
	var networks = make([]string, 0)
	networks = append(networks, "fabricnetwork")
	var ports = make([]string, 0)
	port := nc.PortManager.GetOrdererPort(ordFQDN)
	ports = append(ports, port)

	var orderer Container
	orderer.ContainerName = ordFQDN
	orderer.Extends = extnds
	orderer.Volumns = vols
	orderer.Ports = ports
	if nc.IsMultiMachine() {
		orderer.NetworkMode = "host"
	} else {
		orderer.Networks = networks
	}
	return orderer
}
