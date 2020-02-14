package main

import (
	"bufio"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
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

	resp, err := RestRequest(EPLogin, login, config)
	if err != nil {
		fmt.Println("Err:", err.Error())
		return
	}
	var response loginResponse
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		fmt.Printf("Can't parse response '%s'\n", response)
		return
	}

	if !checkResponse(response.Status, color.New(color.FgHiRed).Sprintf("Failure\n")+"Check credentials and try again") {
		return
	}

	if response.Status == "success" && len(response.Token) > 0 {
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
	return len(config.User.Username) > 0 && len(config.User.SessionToken) > 0
}
