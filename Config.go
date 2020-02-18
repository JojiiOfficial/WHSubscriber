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
	File string `yml:"-"`

	User struct {
		Username     string
		SessionToken string
	}

	Client struct {
		ServerURL   string `default:"https://wh-share.de/"`
		CallbackURL string `default:"https://yourCallbackDomain.de/"`
		IgnoreCert  bool   `default:"false"`
	}

	Webserver struct {
		HTTP  configHTTPstruct
		HTTPS configTLSStruct
	}
}

//Config for HTTPS
type configTLSStruct struct {
	Enabled       bool   `default:"false"`
	ListenAddress string `default:":443"`
	CertFile      string
	KeyFile       string
}

//Config for HTTP
type configHTTPstruct struct {
	Enabled       bool   `default:"false"`
	ListenAddress string `default:":80"`
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

	isDefault, err := configor.SetupConfig(&config, confFile, func(confI interface{}) interface{} {
		conf := confI.(*ConfigStruct)
		conf.Webserver.HTTP = configHTTPstruct{
			Enabled:       true,
			ListenAddress: ":80",
		}
		return conf
	})

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
	config.File = confFile
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
	if !config.Webserver.HTTP.Enabled && !config.Webserver.HTTPS.Enabled {
		fmt.Printf("Error: You need to enable the server first: 'enabled: true' (in the config)")
		return false
	}

	if config.Webserver.HTTP.Enabled {
		if len(config.Webserver.HTTP.ListenAddress) == 0 {
			log.Println("You need to set the HTTP listenaddress in the config")
			return false
		}
	}

	if config.Webserver.HTTPS.Enabled {
		if len(config.Webserver.HTTPS.ListenAddress) == 0 {
			log.Println("You need to set the address in the config")
			return false
		}

		//Check SSL values
		if len(config.Webserver.HTTPS.CertFile) == 0 {
			log.Println("To enable TLS you need to specify a SSL certificate")
			return false
		}
		if len(config.Webserver.HTTPS.KeyFile) == 0 {
			log.Println("To enable TLS you need to specify a SSL private key")
			return false
		}

		//Check SSL files
		if !gaw.FileExists(config.Webserver.HTTPS.CertFile) {
			log.Println("Can't find the SSL certificate. File not found")
			return false
		}
		if !gaw.FileExists(config.Webserver.HTTPS.KeyFile) {
			log.Println("Can't find the SSL key. File not found")
			return false
		}
	}

	return true
}
