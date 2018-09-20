package common

func BuildCLIContainer(dirPath string, otherConatiners []string, nc *NetworkConfig) Container {
	var cli Container
	cli.ContainerName = "cli"
	cli.Image = "hyperledger/fabric-tools:${IMAGE_TAG}"
	cli.TTY = true
	cli.WorkingDir = "/opt/ws"
	vols := make([]string, 0)
	vols = append(vols, "/var/run/:/host/var/run/")
	vols = append(vols, "./:/opt/ws")
	vols = append(vols, "./chaincode/github.com:/opt/gopath/src/github.com")

	cliEnvironment := make([]string, 0)
	cliEnvironment = append(cliEnvironment, "CORE_PEER_TLS_ENABLED=true")
	cliEnvironment = append(cliEnvironment, "GOPATH=/opt/gopath")
	cliEnvironment = append(cliEnvironment, "CORE_LOGGING_LEVEL=DEBUG")
	cliEnvironment = append(cliEnvironment, "CORE_PEER_ID=cli")

	cli.Environment = cliEnvironment
	cli.Volumns = vols

	if nc.IsMultiMachine() {
		cli.NetworkMode = "host"
	} else {
		var networks = make([]string, 0)
		networks = append(networks, "fabricnetwork")
		cli.Networks = networks
		cli.Depends = otherConatiners
	}
	return cli

}
