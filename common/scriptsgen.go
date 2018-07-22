package common

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"

	util "github.com/suddutt1/fabricnetgenv2/util"
)

const _DOWNLOAD_SCRIPTS = `
#!/bin/bash

export VERSION={{.fabricVersion}}
export ARCH=$(echo "$(uname -s|tr '[:upper:]' '[:lower:]'|sed 's/mingw64_nt.*/windows/')-$(uname -m | sed 's/x86_64/amd64/g')" | awk '{print tolower($0)}')
#Set MARCH variable i.e ppc64le,s390x,x86_64,i386
MARCH="x86_64"


: ${CA_TAG:="$MARCH-$VERSION"}
: ${FABRIC_TAG:="$MARCH-$VERSION"}

echo "===> Downloading platform binaries"
curl https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric/hyperledger-fabric/${ARCH}-${VERSION}/hyperledger-fabric-${ARCH}-${VERSION}.tar.gz | tar xz

`
const _GENEATE_ARTIFACTS_TEMPLATE = `
#!/bin/bash -e
export PWD={{ "pwd" | ToCMDString}}

export FABRIC_CFG_PATH=$PWD
export ARCH=$(uname -s)
export CRYPTOGEN=$PWD/bin/cryptogen
export CONFIGTXGEN=$PWD/bin/configtxgen

function generateArtifacts() {
	
	echo " *********** Generating artifacts ************ "
	echo " *********** Deleting old certificates ******* "
	
        rm -rf ./crypto-config
	
        echo " ************ Generating certificates ********* "
	
        $CRYPTOGEN generate --config=$FABRIC_CFG_PATH/crypto-config.yaml
        
        echo " ************ Generating tx files ************ "
	
		$CONFIGTXGEN -profile OrdererGenesis -outputBlock ./genesis.block
		{{range .channels}}{{$chName := .channelName }}{{$channelId:= $chName | ToLower }}
        $CONFIGTXGEN -profile {{print $chName "Channel"}} -outputCreateChannelTx ./{{print $channelId "channel.tx" }} -channelID {{ print $channelId "channel"}}
		{{end}}

}

generateArtifacts 

cd $PWD

`

const _GENEATE_ARTIFACTS_TEMPLATE_WITHCA = `
#!/bin/bash -e
export PWD={{ "pwd" | ToCMDString}}

export FABRIC_CFG_PATH=$PWD
export ARCH=$(uname -s)
export CRYPTOGEN=$PWD/bin/cryptogen
export CONFIGTXGEN=$PWD/bin/configtxgen

function generateArtifacts() {
	
	echo " *********** Generating artifacts ************ "
	echo " *********** Deleting old certificates ******* "
	
        rm -rf ./crypto-config
	
        echo " ************ Generating certificates ********* "
	
        $CRYPTOGEN generate --config=$FABRIC_CFG_PATH/crypto-config.yaml
        
        echo " ************ Generating tx files ************ "
	
		$CONFIGTXGEN -profile OrdererGenesis -outputBlock ./genesis.block
		{{range .channels}}{{$chName := .channelName }}{{$channelId:= $chName | ToLower }}
        $CONFIGTXGEN -profile {{print $chName "Channel"}} -outputCreateChannelTx ./{{print $channelId "channel.tx" }} -channelID {{ print $channelId "channel"}}
		{{end}}

}
function generateDockerComposeFile(){
	OPTS="-i"
	if [ "$ARCH" = "Darwin" ]; then
		OPTS="-it"
	fi
	cp  docker-compose-template.yaml  docker-compose.yaml
	{{ range .orgs}}
	{{$orgName :=.name | ToUpper }}
	cd  crypto-config/peerOrganizations/{{.domain}}/ca
	PRIV_KEY=$(ls *_sk)
	cd ../../../../
	sed $OPTS "s/{{$orgName}}_PRIVATE_KEY/${PRIV_KEY}/g"  docker-compose.yaml
	{{end}}
}
generateArtifacts 
cd $PWD
generateDockerComposeFile
cd $PWD

`
const _CLEAN_UP_SCRIPT = `
#!/bin/bash
echo "Clearing the old artifacts"
rm -f *.yaml
rm -rf crypto-config
rm -f *.block
rm -f *.tx
rm -f generateConfig.sh
rm -f setFabricEnv.sh
rm -f setPeer.sh
rm -f buildandjoinchannel.sh
rm -f *_install.sh
rm -f *_update.sh
rm -f .env
rm -f portmap.json
rm -f downloadbin.sh
rm -f setupChannels.sh
echo "Done!!!"

`
const _SET_ENVIRONMENT = `
#!/bin/bash
export IMAGE_TAG="x86_64-{{.fabricVersion}}"
export TP_IMAGE_TAG="x86_64-{{.thirdPartyVersion}}"

`
const _DOTENV = `
COMPOSE_PROJECT_NAME=bc

`
const _SET_PEER_SCRIPT = `
#!/bin/bash
export ORDERER_CA=/opt/ws/crypto-config/ordererOrganizations/{{.orderers.domain}}/msp/tlscacerts/tlsca.{{.orderers.domain}}-cert.pem
{{$primechannel := (index .channels 0).channelName }}
if [ $# -lt 2 ];then
	echo "Usage : . setPeer.sh {{range .orgs}}{{.name}}|{{end}} <peerid>"
fi
export peerId=$2
{{range .orgs}}
if [[ $1 = "{{.name}}" ]];then
	echo "Setting to organization {{.name}} peer "$peerId
	export CORE_PEER_ADDRESS=$peerId.{{.domain}}:7051
	export CORE_PEER_LOCALMSPID={{.mspID}}
	export CORE_PEER_TLS_CERT_FILE=/opt/ws/crypto-config/peerOrganizations/{{.domain}}/peers/$peerId.{{.domain}}/tls/server.crt
	export CORE_PEER_TLS_KEY_FILE=/opt/ws/crypto-config/peerOrganizations/{{.domain}}/peers/$peerId.{{.domain}}/tls/server.key
	export CORE_PEER_TLS_ROOTCERT_FILE=/opt/ws/crypto-config/peerOrganizations/{{.domain}}/peers/$peerId.{{.domain}}/tls/ca.crt
	export CORE_PEER_MSPCONFIGPATH=/opt/ws/crypto-config/peerOrganizations/{{.domain}}/users/Admin@{{.domain}}/msp
fi
{{end}}

`
const _BUILD_AND_JOIN_CHANNEL_SCRIPT = `
#!/bin/bash -e
{{ $orderer:= .ordererURL}}
{{ $root := . }}
{{range .channels}}
{{ $channelId := print .channelName "channel" | ToLower }}
echo "Building channel for {{print $channelId}}" 
{{$firstOrg := (index .orgs 0) }}
. setPeer.sh {{$firstOrg}} peer0
export CHANNEL_NAME="{{print $channelId }}"
peer channel create -o {{ print $orderer }} -c $CHANNEL_NAME -f ./{{print $channelId ".tx"}} --tls true --cafile $ORDERER_CA -t 10000
{{ range $index,$orgName :=.orgs}}{{$orgConfig :=  index $root $orgName }}
{{ range $i,$peerId:=$orgConfig.peerNames }}
. setPeer.sh {{$orgName}} {{$peerId}}
export CHANNEL_NAME="{{print $channelId }}"
peer channel join -b $CHANNEL_NAME.block
{{end}}{{end}}{{end}}
`

func GenerateDownloadScripts(nc *NetworkConfig, path string) bool {
	dataMapContainer := nc.GetRootConfig()
	isSuccess, outputBytes := util.LoadTemplate(_DOWNLOAD_SCRIPTS, "download", dataMapContainer)
	if !isSuccess {
		fmt.Println("Unable to generate download binary scripts")
		return false
	}
	ioutil.WriteFile(path+"downloadbin.sh", outputBytes, 0777)
	return true
}

func GenerateGenerateArtifactsScript(nc *NetworkConfig, filename string) bool {
	funcMap := template.FuncMap{
		"ToCMDString": ToCMDString,
		"ToLower":     strings.ToLower,
		"ToUpper":     strings.ToUpper,
	}
	templateToUse := _GENEATE_ARTIFACTS_TEMPLATE
	config := nc.GetRootConfig()

	addCA := nc.IsCARequired()
	if addCA == true {
		templateToUse = _GENEATE_ARTIFACTS_TEMPLATE_WITHCA
	}
	tmpl, err := template.New("generateArtifacts").Funcs(funcMap).Parse(templateToUse)
	if err != nil {
		fmt.Printf("Error in reading template %v\n", err)
		return false
	}

	var outputBytes bytes.Buffer
	err = tmpl.Execute(&outputBytes, config)
	if err != nil {
		fmt.Printf("Error in generating the generateArtifacts file %v\n", err)
		return false
	}
	ioutil.WriteFile(filename, outputBytes.Bytes(), 0777)
	return true
}
func GenerateCleanUpScript(nc *NetworkConfig, filePath string) bool {

	isSuccess, scriptBytes := util.LoadTemplate(_CLEAN_UP_SCRIPT, "cleanup", nc.GetRootConfig())
	if !isSuccess {
		fmt.Println("Error in generating cleanup script")
		return false
	}
	ioutil.WriteFile(filePath, scriptBytes, 0777)
	return true
}
func GenerateOtherScripts(nc *NetworkConfig, basePath string) bool {
	isSucess, scriptBytes := util.LoadTemplate(_SET_ENVIRONMENT, "setEnv", nc.GetRootConfig())
	if !isSucess {
		fmt.Println("Unable to generate setFabricEnv.sh file")
		return false
	}
	ioutil.WriteFile(basePath+"/setFabricEnv.sh", scriptBytes, 0777)
	isSucess, scriptBytes = util.LoadTemplate(_DOTENV, "setEnv", nc.GetRootConfig())
	if !isSucess {
		fmt.Println("Unable to generate .env file")
		return false
	}
	ioutil.WriteFile(basePath+"/.env", scriptBytes, 0666)

	return true
}

func GenerateSetPeerScript(nc *NetworkConfig, filename string) bool {
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	tmpl, err := template.New("setPeer").Funcs(funcMap).Parse(_SET_PEER_SCRIPT)
	if err != nil {
		fmt.Printf("Error in reading template %v\n", err)
		return false
	}
	dataMapContainer := nc.GetRootConfig()

	var outputBytes bytes.Buffer
	err = tmpl.Execute(&outputBytes, dataMapContainer)
	if err != nil {
		fmt.Printf("Error in generating the setpeer.sh file %v\n", err)
		return false
	}
	ioutil.WriteFile(filename, outputBytes.Bytes(), 0777)
	return true
}
func GenerateBuildAndJoinChannelScript(nc *NetworkConfig, filename string) bool {
	funcMap := template.FuncMap{
		"ToCMDString": ToCMDString,
		"ToLower":     strings.ToLower,
	}

	tmpl, err := template.New("buildChannel").Funcs(funcMap).Parse(_BUILD_AND_JOIN_CHANNEL_SCRIPT)
	if err != nil {
		fmt.Printf("Error in reading template %v\n", err)
		return false
	}
	dataMapContainer := nc.GetRootConfig()
	channelMap := make(map[string]interface{})

	orgs, _ := dataMapContainer["orgs"].([]interface{})
	for _, org := range orgs {
		orgConfig := util.GetMap(org)
		peerCount := util.GetNumber(orgConfig["peerCount"])
		peerNames := make([]string, 0)
		fmt.Printf(" Peer count %d\n", peerCount)
		for index := 0; index < peerCount; index++ {
			peerNames = append(peerNames, fmt.Sprintf("peer%d", index))
		}
		orgConfig["peerNames"] = peerNames
		orgName := util.GetString(orgConfig["name"])
		channelMap[orgName] = orgConfig
	}
	channelMap["channels"] = dataMapContainer["channels"]
	//Resolve the orderer name
	ordererConfig := util.GetMap(dataMapContainer["orderers"])
	if nc.IsKafkaOrderer() {
		channelMap["ordererURL"] = fmt.Sprintf("%s0.%s:7050", util.GetString(ordererConfig["ordererHostname"]), util.GetString(ordererConfig["domain"]))
	} else {
		channelMap["ordererURL"] = fmt.Sprintf("%s.%s:7050", util.GetString(ordererConfig["ordererHostname"]), util.GetString(ordererConfig["domain"]))
	}
	var outputBytes bytes.Buffer
	err = tmpl.Execute(&outputBytes, channelMap)
	if err != nil {
		fmt.Printf("Error in generating the channel setup script file %v\n", err)
		return false
	}
	ioutil.WriteFile(filename, outputBytes.Bytes(), 0777)
	return true
}

func ToCMDString(input string) string {
	return "`" + input + "`"
}
