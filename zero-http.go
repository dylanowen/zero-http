package main

import (
	"github.com/dylanowen/zero-http/config"
	"github.com/dylanowen/zero-http/server"
	"github.com/spf13/pflag"
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

	var zeroServer = server.NewServer(configuration.Server)

	// start our server
	go func() {
		if err := zeroServer.Start(); err != nil {
			log.Println("Failed to start the server: ", err)
		}

		// tell our main thread to shutdown
		done <- true
	}()

	// listen for a shutdown
	go func() {
		<-interrupt
		log.Println("Shutting Down (ctrl+c to force it)")

		// don't let the server block shutdown
		go func() {
			zeroServer.Stop()
		}()

		// if we get another interrupt force a shutdown
		<-interrupt
		log.Fatalln("Shutdown forced")
	}()

	// wait until our serves completes
	<-done
}

func loadConfig() *config.Configuration {
	// bind all the command line arguments
	pflag.String("port", "", "The port to use")
	pflag.String("certFile", "", "the certFile to use for TLS")
	pflag.String("keyFile", "", "the keyFile to use for TLS")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// load in our configs
	viper.SetConfigName("config")
	viper.AddConfigPath("./.zero-http")
	viper.AddConfigPath("$HOME/.zero-http")

	var configDir = ""
	// check to see if we can find an actual config to load
	if err := viper.MergeInConfig(); err != nil {
		log.Println("Couldn't find a config file to load:", err)
	} else {
		var configFile = viper.ConfigFileUsed()
		configDir = path.Dir(configFile)

		log.Println("Loaded config file:", configFile)
	}

	var rawConfig config.RawConfiguration
	if err := viper.Unmarshal(&rawConfig); err != nil {
		log.Println("Error parsing config file:", err)
	}

	if rawConfig.Debug {
		log.Printf("Config: %+v", rawConfig)
	}

	return config.NewConfiguration(rawConfig, configDir)
}
