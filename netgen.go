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
	genrateCC := false
	flag.BoolVar(&genrateCC, "default-cc-gen", false, "Generate default chaincode[Optional]")
	isSingleMachConfExample := false
	flag.BoolVar(&isSingleMachConfExample, "help-single", false, "Generate example network-config file (Single machine)[Optional]")
	isMultiMachConfExample := false
	flag.BoolVar(&isMultiMachConfExample, "help-multi", false, "Generate example network-config file (Multiple machine)[Optional]")
	isHelp := false
	flag.BoolVar(&isHelp, "help", false, "Prints this help text ")
	flag.Parse()
	if isHelp {
		flag.Usage()
		os.Exit(0)
	}
	if isSingleMachConfExample {
		common.GenerateExampleConfigSingleMachine()
		os.Exit(0)
	}
	if isMultiMachConfExample {
		common.GenerateExampleConfigMultiMachine()
		os.Exit(0)
	}
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
	if genrateCC {
		common.CreateDefaultCC(networkConfig, "./")
	}
	fmt.Println("\nCompleted network configuration generation successfully ")
}
