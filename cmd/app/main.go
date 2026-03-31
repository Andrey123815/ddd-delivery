package main

import (
	"delivery/cmd"
	"delivery/cmd/app/utils"
)

func main() {
	config := utils.GetConfigs()

	compositionRoot := cmd.NewCompositionRoot(config)
	defer compositionRoot.CloseAll()

	utils.StartKafkaProducers(compositionRoot)
	utils.StartKafkaConsumers(compositionRoot)
	utils.StartCrons(compositionRoot)
	utils.StartWebServer(compositionRoot, config.HttpPort)
}
