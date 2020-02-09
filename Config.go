package main

import (
	"log"
	"os"

	"github.com/JojiiOfficial/configor"
)

//ConfigStruct the structure of the configfile
type ConfigStruct struct {
	Server struct {
		Enable  bool `default:"false"`
		Port    int  `default:"443"`
		SSLCert string
		SSLKey  string
	}

	Client struct {
		ServerURL   string `default:"https://wh-share.de/"`
		CallbackURL string `default:"https://yourCallbackDomain.de/"`
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
