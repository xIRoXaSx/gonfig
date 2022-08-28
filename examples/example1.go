package examples

import (
	"log"

	"github.com/xiroxasx/gonfig"
)

type config1 struct {
	BackendAddress string `json:"BackendAddress"`
	Debug          bool   `json:"Debug"`
}

// Example of writing a config1 struct into a file.
func writeConfig1() {
	// The config which we need to store.
	myConfig := config1{
		BackendAddress: "10.0.0.1",
		Debug:          true,
	}

	// Create a new Gonfig instance.
	// The folder containing the configuration will be called "GonfigExample" and will reside in
	// the default root directory to use for user-specific configuration data (https://pkg.go.dev/os#UserConfigDir).
	g, err := gonfig.New("GonfigExample", "config", gonfig.GonfJson, false)
	if err != nil {
		log.Fatal(err)
	}

	// Store the configuration locally.
	err = g.WriteToFile(myConfig)
	if err != nil {
		log.Printf("unable to store config: %v\n", err)
	}
	log.Printf("The config has been saved here: %s\n", g.FullPath())
}

// Example of loading the configuration into a variable.
func loadConfig1() {
	// Create a new Gonfig instance.
	g, err := gonfig.New("GonfigExample", "config", gonfig.GonfJson, false)
	if err != nil {
		log.Fatal(err)
	}

	// Load the configuration from the given path.
	var myConfig config1
	err = g.Load(&myConfig)
	if err != nil {
		log.Fatal(err)
	}
}
