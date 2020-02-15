package main

import (
	"bufio"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	gaw "github.com/JojiiOfficial/GoAw"
	"github.com/JojiiOfficial/configor"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

//LoginCommand login into the server
func LoginCommand(config *ConfigStruct, usernameArg string) {
	inpReader := bufio.NewReader(os.Stdin)
	if isLoggedIn(config) && !*appYes {
		i, e := gaw.ConfirmInput("You are already logged in. Overwrite session? [y/n]> ", inpReader)
		if e == -1 || !i {
			return
		}
	}

	username, pass := credentials(usernameArg)

	// You salty baby
	password := SHA512(pass + SHA512(pass)[:8])

	login := loginRequest{
		Password: password,
		Username: username,
	}

	var response loginResponse
	resp, err := RestRequest(EPLogin, login, &response, config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if resp.Status == ResponseError && resp.HTTPCode == 403 {
		fmt.Println(color.HiRedString("Failure"))
	} else if resp.Status == ResponseSuccess && len(response.Token) > 0 {
		config.User = struct {
			Username     string
			SessionToken string
		}{
			Username:     username,
			SessionToken: response.Token,
		}
		err := configor.Save(config, config.File)
		if err != nil {
			fmt.Println("Error saving config:", err.Error())
			return
		}
		fmt.Println(color.HiGreenString("Success!"), "\nLogged in as", username)
	} else {
		fmt.Println("Unexpected error occured!")
	}
}

func credentials(buser string) (string, string) {
	reader := bufio.NewReader(os.Stdin)
	var username string
	if len(buser) > 0 {
		username = buser
	} else {
		fmt.Print("Enter Username: ")
		username, _ = reader.ReadString('\n')
	}
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalln("Error:", err.Error())
		return "", ""
	}
	return strings.TrimSpace(username), strings.TrimSpace(string(bytePassword))
}

//SHA512 hashes using sha1 algorithm
func SHA512(text string) string {
	algorithm := sha512.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func isLoggedIn(config *ConfigStruct) bool {
	return len(config.User.Username) > 0 && len(config.User.SessionToken) == 64
}

//return true if is logged in
func checkLoggedIn(config *ConfigStruct) bool {
	if isLoggedIn(config) {
		return true
	}
	fmt.Println("You need to be logged in to use this feature")
	return false
}
