package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	gaw "github.com/JojiiOfficial/GoAw"
	"github.com/JojiiOfficial/configService"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

//LoginCommand login into the server
func LoginCommand(config *ConfigStruct, usernameArg string, args ...bool) {
	inpReader := bufio.NewReader(os.Stdin)
	if isLoggedIn(config) && !*appYes && len(args) == 0 {
		i, e := gaw.ConfirmInput("You are already logged in. Overwrite session? [y/n]> ", inpReader)
		if e == -1 || !i {
			return
		}
	}

	username, pass := credentials(usernameArg, false, 0)

	// You salty baby
	password := saltPass(pass)

	login := credentialsRequest{
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
		err := configService.Save(config, config.File)
		if err != nil {
			fmt.Println("Error saving config:", err.Error())
			return
		}
		fmt.Println(color.HiGreenString("Success!"), "\nLogged in as", username)
	} else {
		fmt.Println("Unexpected error occurred!")
	}
}

//RegisterCommand create a new account
func RegisterCommand(config *ConfigStruct) {
	username, pass := credentials("", true, 0)
	if len(username) == 0 || len(pass) == 0 {
		return
	}

	req := credentialsRequest{
		Username: username,
		Password: saltPass(pass),
	}

	resp, err := RestRequest(EPUserCreate, req, nil, config)
	if err != nil {
		fmt.Println("Err", err.Error())
		return
	}

	if resp.Status == ResponseSuccess {
		fmt.Printf("User '%s' created %s!\n", username, color.HiGreenString("successfully"))
		y, _ := gaw.ConfirmInput("Do you want to login to this account? [y/n]> ", bufio.NewReader(os.Stdin))
		if y {
			LoginCommand(config, username, true)
		}
	} else {
		fmt.Println("Error:", resp.Message)
	}
}

func credentials(bUser string, repeat bool, index uint8) (string, string) {
	if index >= 3 {
		return "", ""
	}

	reader := bufio.NewReader(os.Stdin)
	var username string
	if len(bUser) > 0 {
		username = bUser
	} else {
		fmt.Print("Enter Username: ")
		username, _ = reader.ReadString('\n')
	}

	if len(username) > 30 {
		fmt.Println("Username too long!")
		return "", ""
	}

	var pass string
	enterPassMsg := "Enter Password: "
	count := 1

	if repeat {
		count = 2
	}

	for i := 0; i < count; i++ {
		fmt.Print(enterPassMsg)
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalln("Error:", err.Error())
			return "", ""
		}
		fmt.Println()
		lPass := strings.TrimSpace(string(bytePassword))

		if len(lPass) > 80 {
			fmt.Println("Your password is too long!")
			return credentials(username, repeat, index+1)
		}
		if len(lPass) < 7 {
			fmt.Println("Your password must have at least 7 characters!")
			return credentials(username, repeat, index+1)
		}

		if repeat && i == 1 && pass != lPass {
			fmt.Println("Passwords don't match!")
			return credentials(username, repeat, index+1)
		}

		pass = lPass
		enterPassMsg = "Enter Password again: "
	}

	return strings.TrimSpace(username), pass
}

func saltPass(pass string) string {
	return gaw.SHA512(pass + gaw.SHA512(pass)[:10])
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
