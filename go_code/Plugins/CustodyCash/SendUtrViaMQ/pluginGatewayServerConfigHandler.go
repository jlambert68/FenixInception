package SendUtrViaMQ

import (
	"github.com/BurntSushi/toml"
	"log"
)

func processConfigFile(filename string) {
	// If no filename then use predefine filename
	if len(filename) == 0 {
		filename = "pluginDbEngineConfig.toml"
	}

	//	var gatewayConfig TomlConfigStruct
	// Decode toml-file
	if _, err := toml.DecodeFile(filename, &gatewayConfig); err != nil {
		// Error when decoding toml-file
		log.Println("71bb5630-7358-41ba-a01f-b0cb0baade2f")
		log.Println("err: ", err)
		log.Fatalln("Couldn't decode 'toml-file', stopping gateway: ", filename)
	} else {
		// OK when decoding toml-file
		log.Println("5aed5603-c933-4b16-bea1-e6ad0ea9ee75")
		log.Println("'toml-file' was decoded: ", filename)

		// Convert logging-level-text into logrus.logginglever
		log.Println(gatewayConfig)
		//log.Println(gatewayConfig.LoggingLevel.LoggingLevel.String())

	}

}
