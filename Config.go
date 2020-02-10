package main

import (
	"log"
	"os"
	"strings"

	"github.com/JojiiOfficial/configor"
)

//ConfigStruct the structure of the configfile
type ConfigStruct struct {
	Server struct {
		Enable      bool   `default:"false"`
		CallbackURL string `default:"https://yourCallbackDomain.de/"`
		Port        uint16 `default:"443"`
		EnableHTTPS bool
		SSLCert     string
		SSLKey      string
	}

	Client struct {
		ServerURL string `default:"https://wh-share.de/"`
	}
}

var config ConfigStruct

//InitConfig inits the config
//Returns true if system should exit
func InitConfig(confFile string, createMode bool) bool {
	if createMode {
		s, err := os.Stat(confFile)
		if err == nil {
			log.Fatalln("This config already exists!")
			return true
		}
		if s != nil && s.IsDir() {
			log.Fatalln("This name is already taken by a folder")
			return true
		}
		if !strings.HasSuffix(confFile, ".yml") {
			log.Fatalln("The configfile must end with .yml")
			return true
		}
	} else if len(confFile) == 0 {
		confFile = DefaultConfigFile
	}

	isDefault, err := configor.SetupConfig(&config, confFile, configor.NoChange)
	if err != nil {
		log.Fatalln(err.Error())
		return true
	}
	if isDefault {
		log.Println("New config created.")
		log.Println("Exiting")
		return true
	}

	if err = configor.Load(&config, confFile); err != nil {
		log.Fatalln(err.Error())
		return true
	}

	return false
}

//Check check the config file of logical errors
//Returns true on success
func (config *ConfigStruct) Check() bool {
	if config.Server.Enable && len(config.Server.CallbackURL) == 0 {
		log.Println("You need to enter a callbackURL to enable the server")
		return false
	}
	if !inPortValid(config.Server.Port) {
		log.Println("The specified port is invalid")
		return false
	}
	if config.Server.EnableHTTPS {
		if len(config.Server.SSLCert) == 0 {
			log.Println("To enable HTTPS you need to specify a SSL certificate")
			return false
		}
		if len(config.Server.SSLKey) == 0 {
			log.Println("To enable HTTPS you need to specify a SSL private key")
			return false
		}
	}
	if config.Server.Port == 443 && config.Server.EnableHTTPS {
		log.Println("Warning: You shouldn't use HTTP on port 443")
	}

	return true
}
