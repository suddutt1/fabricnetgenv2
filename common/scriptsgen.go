package common

import (
	"fmt"
	"io/ioutil"

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

func GenerateDownloadScripts(nc *NetworkConfig, path string) bool {
	nc.DetermineImageVersions()
	dataMapContainer := nc.GetRootConfig()
	isSuccess, outputBytes := util.LoadTemplate(_DOWNLOAD_SCRIPTS, "download", dataMapContainer)
	if !isSuccess {
		fmt.Println("Unable to generate download binary scripts")
		return false
	}
	ioutil.WriteFile(path+"downloadbin.sh", outputBytes, 0777)
	return true
}
