package FenixTestExecutionServer

import (
	"github.com/BurntSushi/toml"
	"log"
)

func processConfigFile(filename string) {
	// If no filename then use predefine filename
	if len(filename) == 0 {
		filename = "gatewayConfig.toml"
	}

	//	var gatewayConfig TomlConfigStruct
	// Decode toml-file
	if _, err := toml.DecodeFile(filename, &gatewayConfig); err != nil {
		// Error when decoding toml-file
		log.Println("2243085a-feee-4ae7-8ccf-03f69c0704a4")
		log.Println("err: ", err)
		log.Fatalln("Couldn't decode 'toml-file', stopping gateway: ", filename)
	} else {
		// OK when decoding toml-file
		log.Println("2d8a80d96-6c58-4f26-a341-71b15a760569")
		log.Println("'toml-file' was decoded: ", filename)

		// Convert logging-level-text into logrus.logginglever
		log.Println(gatewayConfig)
		//log.Println(gatewayConfig.LoggingLevel.LoggingLevel.String())

	}

}
