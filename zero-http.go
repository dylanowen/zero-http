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

	var servers = []*server.ZeroServer{
		server.NewServer(configuration.Http, configuration.ConfigDir),
	}

	if configuration.Https != nil && configuration.Https.Port > 0 {
		servers = append(servers, server.NewServer(configuration.Https, configuration.ConfigDir))
	}

	// start our servers
	log.Println("Starting Up")
	for i := 0; i < len(servers); i++ {
		var s = servers[i]

		go func() {
			if err := s.Start(); err != nil {
				log.Println("Failed to start the server: ", err)
			}

			// tell our main thread to shutdown
			done <- true
		}()
	}

	// listen for a shutdown
	go func() {
		<-interrupt
		log.Println("Shutting Down (ctrl+c to force it)")

		for i := 0; i < len(servers); i++ {
			var s = servers[i]
			// don't let the server block shutdown
			go func() {
				s.Stop()
			}()
		}

		// if we get another interrupt force a shutdown
		<-interrupt
		log.Fatalln("Shutdown forced")
	}()

	// wait until our servers complete
	for i := 0; i < len(servers); i++ {
		<-done
	}
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
