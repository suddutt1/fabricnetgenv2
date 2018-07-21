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
rm -f setenv.sh
rm -f setpeer.sh
rm -f buildandjoinchannel.sh
rm -f *_install.sh
rm -f *_update.sh
rm -f .env
rm -f portmap.json
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
func ToCMDString(input string) string {
	return "`" + input + "`"
}
