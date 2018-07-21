package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	common "github.com/suddutt1/fabricnetgenv2/common"
)

func main() {
	fmt.Println("Starting fabric network generator V2.0")
	configFile := ""
	flag.StringVar(&configFile, "net-config", "", "Please provide the network config json file")
	flag.Parse()
	if len(configFile) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	fmt.Println("Using config file ", configFile)
	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println("Unable to read the config file ", configFile)
		os.Exit(2)
	}
	networkConfig := new(common.NetworkConfig)
	err = json.Unmarshal(configBytes, &networkConfig)
	if err != nil {
		fmt.Println("Unable to parse the config file ", configFile)
		os.Exit(2)
	}
	if !common.GenerateConfigForSingleMachine(networkConfig, "./") {
		fmt.Println("Unable to generate network configuration")
		os.Exit(3)
	}
	fmt.Println("Completed network configuration generation successfully ")
}
