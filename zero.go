package main

import (
	"github.com/dylanowen/zero-http/config"
	"github.com/dylanowen/zero-http/server"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
)

func main() {
	var configuration = loadConfig()

	var interrupt = make(chan os.Signal)
	var done = make(chan bool)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var zeroServer = server.NewZeroServer(&configuration.Server)

	// start our action
	go func() {
		log.Println("Starting Up")

		if err := zeroServer.Start(); err != nil {
			log.Fatalln("Failed to start the server: ", err)
		}

		// tell our main thread to shutdown
		done <- true
	}()

	// listen for a shutdown
	go func() {
		<-interrupt
		log.Println("Shutting Down (ctrl+c to force it)")

		// don't let the action block shutdown
		go func() {
			zeroServer.Stop()
		}()

		// if we get another interrupt force a shutdown
		<-interrupt
		log.Fatalln("Shutdown forced")
	}()

	// wait until our action completes
	<-done
}

func loadConfig() *config.Configuration {
	// bind all the possible command line arguments
	config.ParseCommandLine(viper.GetViper())

	// set our defaults
	config.SetDefault(viper.GetViper())

	// load in our configs
	viper.SetConfigName("config")
	viper.AddConfigPath("./.zero-http")
	viper.AddConfigPath("$HOME/.zero-http")

	// merge the actual config on top of the defaults
	if err := viper.MergeInConfig(); err != nil {
		log.Println("Couldn't find a config file to load:", err)
	}

	var rawConfig config.RawConfiguration
	if err := viper.Unmarshal(&rawConfig); err != nil {
		log.Println("Error parsing config file:", err)
	}

	var configFile = viper.ConfigFileUsed()
	var configDir = path.Dir(configFile)

	log.Println("Loaded config file:", configFile)

	if rawConfig.Debug {
		log.Printf("Config: %+v", rawConfig)
	}

	return &config.Configuration{
		RawConfiguration: rawConfig,
		ConfigDir:        configDir,
	}
}
