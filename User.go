package main

import (
	"bufio"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	gaw "github.com/JojiiOfficial/GoAw"
	"github.com/JojiiOfficial/configor"
	"github.com/fatih/color"
)

//LoginCommand login into the server
func LoginCommand(config *ConfigStruct, usernameArg string) {
	inpReader := bufio.NewReader(os.Stdin)
	if isLoggedIn(config) {
		i, e := gaw.ConfirmInput("You are already logged in. Overwrite session? [y/n]> ", inpReader)
		if e == -1 || !i {
			return
		}
	}

	var username, password string
	if len(usernameArg) > 0 {
		username = usernameArg
	} else {
		i, text := gaw.WaitForMessage("Username: ", inpReader)
		if i != 1 {
			fmt.Println("Abort")
		}
		username = text
	}
	i, text := gaw.WaitForMessage("Password: ", inpReader)
	if i != 1 {
		fmt.Println("Abort")
	}
	// You salty baby
	password = SHA512(text + SHA512(text)[:8])
	fmt.Println(password)

	login := loginRequest{
		Password: password,
		Username: username,
	}

	resp, err := RestRequest(EPLogin, login, config)
	if err != nil {
		fmt.Println("Err:", err.Error())
	}
	var response loginResponse
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		fmt.Printf("Can't parse response '%s'\n", response)
		return
	}
	if response.Status == "Error" {
		fmt.Println("Error logging in. Check credentials and try again")
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
		fmt.Println(color.HiGreenString("Success!"), "Logged in as", username)
	} else {
		fmt.Println("Unexpected error occured!")
	}
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
