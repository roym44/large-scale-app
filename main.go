package main

import (
	"log"
	"os"

	"github.com/TAULargeScaleWorkshop/RLAD/config"
	"github.com/TAULargeScaleWorkshop/RLAD/utils"
	"gopkg.in/yaml.v2"
)

func main() {
	// read configuration file from command line argument
	if len(os.Args) != 2 {
		utils.Logger.Fatal("Expecting exactly one configuration file")
		os.Exit(1)
	}
	configFile := os.Args[1]
	configData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
		os.Exit(2)
	}

	var config config.ConfigBase
	err = yaml.Unmarshal(configData, &config) // parses YAML
	if err != nil {
		log.Fatalf("error unmarshaling data: %v", err)
		os.Exit(3)
	}

	switch config.Type {
	case "TestService":
		utils.Logger.Printf("Loading service type: %v\n", config.Type)
	default:
		utils.Logger.Fatalf("Unknown configuration type: %v", config.Type)
		os.Exit(4)
	}
}
