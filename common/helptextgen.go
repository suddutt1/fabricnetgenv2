package common

import (
	"encoding/json"
	"fmt"
)

func GenerateExampleConfigSingleMachine() {
	fmt.Printf("\n%s\n", _EXAMPLE_CONFIG_SIGLE_MACH)
}
func GenerateExampleConfigMultiMachine() {
	fmt.Printf("\n%s\n", _EXAMPLE_CONFIG_MULTI_MACH)
}
func prettyPrint(str string) string {
	var genericObj interface{}
	json.Unmarshal([]byte(str), &genericObj)
	prettyBytes, _ := json.MarshalIndent(genericObj, "", " ")
	return string(prettyBytes)
}

const _EXAMPLE_CONFIG_MULTI_MACH = `
{
    "description":"Product provenance and traceability sample network",
    "multiMachine":"true",
	"fabricVersion":"1.1.0",
    "orderers":{
        "name" :"Orderer","mspID":"OrdererMSP","domain":"orderer.net","ordererHostname":"orderer","SANS":"localhost","type":"solo"
    },
    "addCA":"false",
    "orgs":[
        { 
            "name" :"Manufacturer",
            "domain":"manuf.net",
            "mspID":"ManufacturerMSP",
            "SANS":"localhost",
            "peerCount":1,
            "userCount":1
        },
        { 
            "name" :"Distributer",
            "domain":"distributer.net",
            "mspID":"DistributerMSP",
            "SANS":"localhost",
            "peerCount":1,
            "userCount":1
        },
        { 
            "name" :"Retailer",
            "domain":"retailer.com",
            "mspID":"RetailerMSP",
            "SANS":"localhost",
            "peerCount":1,
            "userCount":1
        },
        { 
            "name" :"Consumer",
            "domain":"consumerportal.net",
            "mspID":"ConsumerMSP",
            "SANS":"localhost",
            "peerCount":1,
            "userCount":1
        }
        ],
    "consortium":"SupplyChainConsortium",
    "channels" :[
                    {"channelName":"prodtracking","orgs":["Manufacturer","Distributer","Retailer","Consumer"] },
                    {"channelName":"settlement","orgs":["Manufacturer","Distributer","Retailer"] }
                ],
    "chaincodes":[{"channelName":"prodtracking","ccid":"prodtracer","version":"1.0","src":"github.com/prodtracer","participants":["Manufacturer","Distributer","Retailer","Consumer"]},
        {"channelName":"settlement","ccid":"bsmgmt","version":"1.0","src":"github.com/bsmgmt","participants":["Manufacturer","Distributer","Retailer"]}
]            
                
}

`
const _EXAMPLE_CONFIG_SIGLE_MACH = `
{
    "description":"Product provenance and traceability sample network",
    "fabricVersion":"1.1.0",
    "orderers":{
        "name" :"Orderer","mspID":"OrdererMSP","domain":"orderer.net","ordererHostname":"orderer","SANS":"localhost","type":"solo"
    },
    "addCA":"false",
    "startPort":"2000",
    "orgs":[
        { 
            "name" :"Manufacturer",
            "domain":"manuf.net",
            "mspID":"ManufacturerMSP",
            "SANS":"localhost",
            "peerCount":1,
            "userCount":1
        },
        { 
            "name" :"Distributer",
            "domain":"distributer.net",
            "mspID":"DistributerMSP",
            "SANS":"localhost",
            "peerCount":1,
            "userCount":1
        },
        { 
            "name" :"Retailer",
            "domain":"retailer.com",
            "mspID":"RetailerMSP",
            "SANS":"localhost",
            "peerCount":1,
            "userCount":1
        },
        { 
            "name" :"Consumer",
            "domain":"consumerportal.net",
            "mspID":"ConsumerMSP",
            "SANS":"localhost",
            "peerCount":1,
            "userCount":1
        }
        ],
    "consortium":"SupplyChainConsortium",
    "channels" :[
                    {"channelName":"prodtracking","orgs":["Manufacturer","Distributer","Retailer","Consumer"] },
                    {"channelName":"settlement","orgs":["Manufacturer","Distributer","Retailer"] }
                ],
    "chaincodes":[{"channelName":"prodtracking","ccid":"prodtracer","version":"1.0","src":"github.com/prodtracer","participants":["Manufacturer","Distributer","Retailer","Consumer"]},
        {"channelName":"settlement","ccid":"bsmgmt","version":"1.0","src":"github.com/bsmgmt","participants":["Manufacturer","Distributer","Retailer"]}
]            
                
}

`
