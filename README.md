# Gonfig
Configs, simplified 🎉.

## Description
Originated from various projects, Gonfig is a minimalistic package which enables you to create configurations with ease.  
With simplicity in mind, Gonfig will take care of your project's configuration files, so you can focus on more relevant 
aspects of your project.

## Usage
The usage is pretty straight forward.  
Install the latest package via `go get github.com/xIRoXaSx/gonfig` and use as shown in the example.

### Example - Minimalistic approach
In this example we're going to store some configuration inside the users default root configuration directory.  
Please consult the examples directory for further details. 
```go
type config struct {
    BackendAddress string `json:"BackendAddress"`
    Debug          bool   `json:"Debug"`
}

// Example of writing a config struct into a file and reading it.
func writeConfig() {
    // The config which we need to store.
    myConfig := config{
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

    // Load the configuration from the given path.
    var myConfig config
    err = g.Load(&myConfig) // Notice the pointer!
    if err != nil {
        log.Fatal(err)
    }
}
```
