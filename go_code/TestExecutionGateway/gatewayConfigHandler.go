package TestExecutionGateway

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
		log.Println("045bcaee-eefe-419d-8c3d-51ca5adac22e")
		log.Println("err: ", err)
		log.Fatalln("Couldn't decode 'toml-file', stopping gateway: ", filename)
	} else {
		// OK when decoding toml-file
		log.Println("7706838d-72d8-4311-ba39-4efebef75dcf")
		log.Println("'toml-file' was decoded: ", filename)

		// Convert logging-level-text into logrus.logginglever
		log.Println(gatewayConfig)
		//log.Println(gatewayConfig.LoggingLevel.LoggingLevel.String())

	}

}

// *****************************************************************
// Convert logging-level-text into logrus.logginglever
//
/*
func convertLoggingLevel() {
	switch gatewayConfig.LoggingLevel.LoggingLevel {
	case "Debug":


	case "Info":

	case "Warning":

	default:

	}

*/

/*

func (d *logrus.Level) UnmarshalText(text []byte) error {
	var err error
	d, err = time.ParseDuration(string(text))
	return err
}
*/
