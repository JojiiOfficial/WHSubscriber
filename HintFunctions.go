package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func hintListCurDir() []string {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println("Error listing cur dir:", err.Error())
		return []string{}
	}
	var s []string

	for _, e := range files {
		if e.IsDir() {
			s = append(s, e.Name())
		}
	}
	return s
}

func hintRandomNames() []string {
	c := 5
	names := make([]string, c)
	for i := 0; i < c; i++ {
		names[i] = getRandomName()
	}
	return names
}

func getDBunparsed() string {
	db := os.Getenv(getEnVar(EnVarDatabaseFile))
	if len(db) == 0 {
		db = getDefaultDBFile()
	}
	return db
}

func hintListActionIDs() []string {
	db, err := connectDB(getDBunparsed())
	if err != nil {
		fmt.Println(err.Error())
	}
	dat, err := getActionIDs(db, 9)
	if err != nil {
		log.Fatalln(err.Error())
		return []string{}
	}
	return dat
}

func hintListActions() []string {
	var ret []string
	db, err := connectDB(getDBunparsed())
	if err != nil {
		fmt.Println(err.Error())
	}
	actions, err := getActions(db)
	if err != nil {
		log.Fatalln(err.Error())
		return []string{}
	}
	for _, action := range actions {
		ret = append(ret, action.Name)
	}
	return ret
}

func hintSubscriptionsNoNa() []string {
	return hintSubscriptionsB(false)
}
func hintSubscriptions() []string {
	return hintSubscriptionsB(true)
}
func hintSubscriptionsB(add bool) []string {
	db, err := connectDB(getDBunparsed())
	if err != nil {
		fmt.Println(err.Error())
	}
	dat, err := getSubscriptionsHumanized(db, 9)
	if err != nil {
		log.Fatalln(err.Error())
		return []string{}
	}
	ap := 0
	if add {
		ap = 1
	}
	data := make([]string, len(dat)+ap)
	for i, v := range dat {
		data[i] = v
	}
	if add {
		data[len(dat)] = "na"
	}
	return data
}

func hintAvailableActions() []string {
	s := []string{}
	for action := range Modes {
		s = append(s, action)
	}
	return s
}

func hintAvailableActionsForSource() []string {
	s := []string{}
	for action := range Modes {
		if action == "script" {
			action = "custom"
		}
		s = append(s, action)
	}
	return s
}
