package PluginKeyValueDBStore

import (
	"github.com/BurntSushi/toml"
	"log"
)

func processConfigFile(filename string) {
	// If no filename then use predefine filename
	if len(filename) == 0 {
		filename = "pluginDbEngineConfig.toml"
	}

	//	var keyValueStoreConfig TomlConfigStruct
	// Decode toml-file
	if _, err := toml.DecodeFile(filename, &keyValueStoreConfig); err != nil {
		// Error when decoding toml-file
		log.Println("cbd5e603-ae31-437b-9139-eaeece355b95")
		log.Println("err: ", err)
		log.Fatalln("Couldn't decode 'toml-file': ", filename, ". Stopping KeyValueStoreDB")
	} else {
		// OK when decoding toml-file
		log.Println("483dea7d-85d3-4009-93ed-be418a5048e1")
		log.Println("'toml-file':", filename, " was decoded")

		// Convert logging-level-text into logrus.logginglever
		log.Println(keyValueStoreConfig)
		//log.Println(keyValueStoreConfig.LoggingLevel.LoggingLevel.String())

	}

}
