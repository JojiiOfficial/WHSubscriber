package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	gaw "github.com/JojiiOfficial/GoAw"
	"github.com/JojiiOfficial/configor"
)

//ConfigStruct the structure of the configfile
type ConfigStruct struct {
	Client struct {
		ServerURL   string `default:"https://wh-share.de/"`
		CallbackURL string `default:"https://yourCallbackDomain.de/"`
	}

	Server struct {
		Enable        bool   `default:"false"`
		ListenAddress string `default:":8499"`
		UseTLS        bool
		SSLCert       string
		SSLKey        string
	}
}

func getDefaultConfig() string {
	return path.Join(getDataPath(), DefaultConfigFile)
}

//InitConfig inits the config
//Returns true if system should exit
func InitConfig(confFile string, createMode bool) (*ConfigStruct, bool) {
	var config ConfigStruct
	if len(confFile) == 0 {
		confFile = getDefaultConfig()
	}
	if createMode {
		s, err := os.Stat(confFile)
		if err == nil {
			log.Fatalln("This config already exists!")
			return nil, true
		}
		if s != nil && s.IsDir() {
			log.Fatalln("This name is already taken by a folder")
			return nil, true
		}
		if !strings.HasSuffix(confFile, ".yml") {
			log.Fatalln("The configfile must end with .yml")
			return nil, true
		}
	}

	isDefault, err := configor.SetupConfig(&config, confFile, configor.NoChange)
	if err != nil {
		log.Fatalln(err.Error())
		return nil, true
	}
	if isDefault {
		log.Println("New config created.")
		if createMode {
			log.Println("Exiting")
			return nil, true
		}
	}

	if err = configor.Load(&config, confFile); err != nil {
		log.Fatalln(err.Error())
		return nil, true
	}

	return &config, false
}

//Check check the config file of logical errors
func (config *ConfigStruct) Check() bool {
	return true
}

//CheckServer check the config for the server of logical errors
//Returns true on success
func (config *ConfigStruct) CheckServer() bool {
	//Validate server configuration if enabled
	if !config.Server.Enable {
		fmt.Printf("Error: You need to enable the server first: 'enabled: true' (in the config)")
		return false
	}

	if len(config.Server.ListenAddress) == 0 {
		log.Println("You need to set the address in the config")
		return false
	}

	if config.Server.UseTLS {
		//Check SSL values
		if len(config.Server.SSLCert) == 0 {
			log.Println("To enable TLS you need to specify a SSL certificate")
			return false
		}
		if len(config.Server.SSLKey) == 0 {
			log.Println("To enable TLS you need to specify a SSL private key")
			return false
		}

		//Check SSL files
		if !gaw.FileExists(config.Server.SSLCert) {
			log.Println("Can't find the SSL certificate. File not found")
			return false
		}
		if !gaw.FileExists(config.Server.SSLKey) {
			log.Println("Can't find the SSL key. File not found")
			return false
		}
	}

	u, err := gaw.ParseURL(config.Client.CallbackURL)
	if err != nil {
		log.Println("Can't parse CallbackURL:", err.Error())
		return false
	}
	if len(u.Path) != 0 && u.Path != "/" {
		log.Println("You can't specify a path in the CallbackURL. Use a reverseproxy to change the path")
		return false
	}

	return true
}
