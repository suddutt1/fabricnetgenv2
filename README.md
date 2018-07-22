##Hyperledger Fabric Network Generator V2

### Proposed features
1. Single Machine config
2. Multiple Machine config
3. K8S based config

Each of the above mentioned configuration will generate the following
1. Required configTx, crypto-config.yaml, docker yamls/K8S yamls
2. Required shell scripts to expidite the installation process
3. README.md file with step by step execution process

### Implemented features
1. Single machine config with solo orderer without CA
2. Contineous port allocation for (Single machine)
3. Scipts to automate following ( Single machine)
    * Download binary
    * Generate configuration files out of configTx file
    * Channel build and join scripts
    * Chaincode build and update sctipts
4. Generation of a default chaincode for example

### TODO:
1. Introduce CA ( Single machine)
2. Introduce Kafka ordering ( Single machine)
3. Multi machine configuration ( Features to be added TBD)
4. K8S configuration    