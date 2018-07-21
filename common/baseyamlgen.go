package common

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func GenerateBaseYAML(nc *NetworkConfig, basePath string) bool {
	var peerbase Container
	peerEnvironment := make([]string, 0)
	peerEnvironment = append(peerEnvironment, "CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock")
	peerEnvironment = append(peerEnvironment, "CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=bc_fabricnetwork")
	peerEnvironment = append(peerEnvironment, "CORE_LOGGING_LEVEL=DEBUG")
	peerEnvironment = append(peerEnvironment, "CORE_PEER_TLS_ENABLED=true")
	peerEnvironment = append(peerEnvironment, "CORE_PEER_ENDORSER_ENABLED=true")
	peerEnvironment = append(peerEnvironment, "CORE_PEER_GOSSIP_USELEADERELECTION=true")
	peerEnvironment = append(peerEnvironment, "CORE_PEER_GOSSIP_ORGLEADER=false")
	peerEnvironment = append(peerEnvironment, "CORE_PEER_PROFILE_ENABLED=true")
	peerEnvironment = append(peerEnvironment, "CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt")
	peerEnvironment = append(peerEnvironment, "CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key")
	peerEnvironment = append(peerEnvironment, "CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt")

	peerbase.Image = "hyperledger/fabric-peer:${IMAGE_TAG}"
	peerbase.Environment = peerEnvironment
	peerbase.WorkingDir = "/opt/gopath/src/github.com/hyperledger/fabric/peer"
	peerbase.Command = "peer node start"
	config := make(map[string]interface{})
	config["peer"] = peerbase

	var ordererBase Container
	ordererMSP := nc.GetOrdererMSP()
	ordererEnvironment := make([]string, 0)
	ordererEnvironment = append(ordererEnvironment, "ORDERER_GENERAL_LOGLEVEL=debug")
	ordererEnvironment = append(ordererEnvironment, "ORDERER_GENERAL_LISTENADDRESS=0.0.0.0")
	ordererEnvironment = append(ordererEnvironment, "ORDERER_GENERAL_GENESISMETHOD=file")
	ordererEnvironment = append(ordererEnvironment, "ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/genesis.block")
	ordererEnvironment = append(ordererEnvironment, "ORDERER_GENERAL_LOCALMSPID="+ordererMSP)
	ordererEnvironment = append(ordererEnvironment, "ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp")
	ordererEnvironment = append(ordererEnvironment, "ORDERER_GENERAL_TLS_ENABLED=true")
	ordererEnvironment = append(ordererEnvironment, "ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key")
	ordererEnvironment = append(ordererEnvironment, "ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt")
	ordererEnvironment = append(ordererEnvironment, "ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]")
	ordererEnvironment = append(ordererEnvironment, "ORDERER_KAFKA_RETRY_SHORTINTERVAL=1s")
	ordererEnvironment = append(ordererEnvironment, "ORDERER_KAFKA_RETRY_SHORTTOTAL=30s")
	ordererEnvironment = append(ordererEnvironment, "ORDERER_KAFKA_VERBOSE=true")

	ordererBase.Image = "hyperledger/fabric-orderer:${IMAGE_TAG}"
	ordererBase.Environment = ordererEnvironment
	ordererBase.WorkingDir = "/opt/gopath/src/github.com/hyperledger/fabric"
	ordererBase.Command = "orderer"
	config["orderer"] = ordererBase

	var couchDB Container
	couchDB.Image = "hyperledger/fabric-couchdb:${TP_IMAGE_TAG}"
	config["couchdb"] = couchDB
	addCA := nc.IsCARequired()
	if addCA == true {
		var ca Container
		ca.Image = "hyperledger/fabric-ca:${IMAGE_TAG}"
		caEnvironment := make([]string, 0)
		caEnvironment = append(caEnvironment, "FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server")
		caEnvironment = append(caEnvironment, "FABRIC_CA_SERVER_TLS_ENABLED=true")
		ca.Environment = caEnvironment
		ca.Command = "sh -c 'fabric-ca-server start -b admin:adminpw -d'"
		config["ca"] = ca
	}
	var zookeeper Container
	zookeeper.Image = "hyperledger/fabric-zookeeper:${TP_IMAGE_TAG}"
	zookeeper.Restart = "always"
	ports := make([]string, 0)
	ports = append(ports, "2181")
	ports = append(ports, "2888")
	ports = append(ports, "3888")
	zookeeper.Ports = ports
	config["zookeeper"] = zookeeper
	var kfka Container
	kfka.Image = "hyperledger/fabric-kafka:${TP_IMAGE_TAG}"
	kfka.Restart = "always"
	kfkaEnv := make([]string, 0)
	kfkaEnv = append(kfkaEnv, "KAFKA_MESSAGE_MAX_BYTES=103809024")
	kfkaEnv = append(kfkaEnv, "KAFKA_REPLICA_FETCH_MAX_BYTES=103809024")
	kfkaEnv = append(kfkaEnv, "KAFKA_UNCLEAN_LEADER_ELECTION_ENABLE=false")
	kfka.Environment = kfkaEnv
	kports := make([]string, 0)
	kports = append(kports, "9092")
	kfka.Ports = kports
	config["kafka"] = kfka

	var serviceConfig ServiceConfig
	serviceConfig.Version = "2"
	serviceConfig.Services = config
	outBytes, _ := yaml.Marshal(serviceConfig)
	ioutil.WriteFile(basePath+"base.yaml", outBytes, 0666)
	return true
}
