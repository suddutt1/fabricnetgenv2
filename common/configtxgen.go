package common

import (
	"fmt"
	"io/ioutil"

	util "github.com/suddutt1/fabricnetgenv2/util"
)

const _CONFIG_TX_TEMPLATE = `
Profiles:

    OrdererGenesis:
        Orderer:
            <<: *OrdererDefaults
            Organizations:
                - *OrdererOrg
        Consortiums:
          {{.consortium}}:
             Organizations:
                {{ range .orgs}}- *{{ .name}}Org
                {{end}}
    {{ $x :=.consortium}}
    {{range .channels}}
    {{.channelName}}Channel:
        Consortium: {{$x}}
        Application:
            <<: *ApplicationDefaults
            Organizations:
                {{range $index,$var := .orgs}}- *{{$var}}Org
                {{end}}
    {{end}} 
Organizations:
    - &OrdererOrg
        Name: {{index .orderers "mspID" }}
        ID: {{index .orderers "mspID" }}
        MSPDir: crypto-config/ordererOrganizations/{{ index .orderers "domain" }}/msp
    {{range .orgs}}
    - &{{ .name}}Org
        Name: {{.mspID}}
        ID: {{.mspID}}
        MSPDir: crypto-config/peerOrganizations/{{ .domain  }}/msp
        AnchorPeers:
          - Host: peer0.{{.domain}}
            Port: 7051
        {{ end }}
{{ if  and (eq .orderers.type "kafka")  (  .orderers.haCount ) }}
Orderer: &OrdererDefaults
        OrdererType: kafka
        Addresses:{{ range .ordererFDQNList }}
          - {{.}}:7050{{end}}
        BatchTimeout: 2s
        BatchSize:
          MaxMessageCount: 10
          AbsoluteMaxBytes: 98 MB
          PreferredMaxBytes: 512 KB
        Kafka:
          Brokers:
            - kafka0:9092
            - kafka1:9092
            - kafka2:9092
            - kafka3:9092
        Organizations:
{{else}}
Orderer: &OrdererDefaults
        OrdererType: solo
        Addresses:
          - {{index .orderers "ordererHostname" }}.{{index .orderers "domain"}}:7050
        BatchTimeout: 2s
        BatchSize:
          MaxMessageCount: 10
          AbsoluteMaxBytes: 98 MB
          PreferredMaxBytes: 512 KB
        Kafka:
          Brokers:
            - 127.0.0.1:9092
        Organizations:

{{end}}    
Application: &ApplicationDefaults
    Organizations:
`

func GenerateConfigTxGen(nc *NetworkConfig, filename string) bool {

	dataMapContainer := nc.GetRootConfig()
	ordererConfig := util.GetMap(dataMapContainer["orderers"])
	if util.IfEntryExistsInMap(ordererConfig, "type") && util.IfEntryExistsInMap(ordererConfig, "haCount") {
		if util.GetString(ordererConfig["type"]) == "kafka" {
			hostName := util.GetString(ordererConfig["ordererHostname"])
			domainName := util.GetString(ordererConfig["domain"])
			listOfOrderers := make([]string, 0)
			for index := 0; index < util.GetNumber(ordererConfig["haCount"]); index++ {
				listOfOrderers = append(listOfOrderers, fmt.Sprintf("%s%d.%s", hostName, index, domainName))
			}
			dataMapContainer["ordererFDQNList"] = listOfOrderers
		}
	}
	isSucess, dataBytes := util.LoadTemplate(_CONFIG_TX_TEMPLATE, "configTx", dataMapContainer)
	if !isSucess {
		fmt.Printf("Error in generating configtx.yaml file\n")
		return false
	}
	err := ioutil.WriteFile(filename, dataBytes, 0666)
	if err != nil {
		fmt.Printf("Error in generating configTx file %v\n", err)
		return false
	}
	return true
}
