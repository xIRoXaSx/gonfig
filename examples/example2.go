package examples

import (
	"log"
	"time"

	"github.com/xiroxasx/gonfig"
)

type config2 struct {
	CreatedAt time.Time `yaml:"Debug"`
}

const configPath = "/SomeWhere/On/Your/Machine/config.yaml"

// Example of writing a config2 struct into a file in a location other than the default configuration directory.
func writeConfig2() {
	// The config which we need to store.
	myConfig := config2{
		CreatedAt: time.Now(),
	}

	// Create a new Gonfig instance.
	// The yaml configuration will now be written to the given absolute path.
	g, err := gonfig.NewWithPath(configPath, gonfig.GonfYaml, false)
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
func loadConfig2() {
	// Create a new Gonfig instance.
	g, err := gonfig.NewWithPath(configPath, gonfig.GonfYaml, false)
	if err != nil {
		log.Fatal(err)
	}

	// Load the configuration from the given path.
	var myConfig config2
	err = g.Load(&myConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("config's full path", g.FullPath())
	log.Println("config's directory", g.Dir())
	log.Println("config's type", g.Type())
}
