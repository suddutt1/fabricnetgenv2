package main

import (
	"fmt"

	common "github.com/suddutt1/fabricnetgenv2/common"
)

func main() {
	fmt.Println("Starting fabric network generator V2.0")
	common.GenerateOrdererConfigK8S()
}
