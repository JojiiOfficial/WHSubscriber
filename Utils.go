package main

import (
	"errors"
	"log"
	"net"
	"net/url"
	"os"
	"os/user"
	"path"
	"strings"

	gaw "github.com/JojiiOfficial/GoAw"
)

func getDataPath() string {
	path := path.Join(gaw.GetHome(), DataDir)
	s, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, 0770)
		if err != nil {
			log.Fatalln(err.Error())
		}
	} else if s != nil && !s.IsDir() {
		log.Fatalln("DataPath-name already taken by a file!")
	}
	return path
}

func mapKeyByValue(val uint8, m map[string]uint8) string {
	for k, v := range m {
		if v == val {
			return k
		}
	}
	return ""
}

//Looks if ip is assigned to host
func matchIPHost(ip, host string) (bool, error) {
	u, err := url.Parse(host)
	if err == nil && len(u.Hostname()) > 0 {
		host = u.Hostname()
	}

	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		return false, errors.New("NoIP")
	}
	ips, err := net.LookupHost(host)
	if err != nil {
		return false, err
	}

	for _, ipa := range ips {
		if ipa == ip {
			return true, nil
		}
	}
	return false, nil
}

func getUsername(custUser ...string) string {
	if len(custUser) > 0 && len(custUser[0]) > 0 {
		return custUser[0]
	}

	user, err := user.Current()
	if err == nil {
		return user.Username
	}
	return ""
}

func formatBashEnVars(enVars []string) string {
	envStr := strings.Join(enVars, "; ")
	if len(enVars) > 0 {
		envStr += ";"
	}
	return envStr
}
