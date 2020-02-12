package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func hintListCurrDir() []string {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println("Error listing curr dir:", err.Error())
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

func hintSubscriptions() []string {
	db, err := connectDB(getDBunparsed())
	if err != nil {
		fmt.Println(err.Error())
	}
	dat, err := getHooksHumanized(db, 9)
	if err != nil {
		log.Fatalln(err.Error())
		return []string{}
	}
	return dat
}

func hintAvailableActions() []string {
	s := []string{}
	for action := range Actions {
		s = append(s, action)
	}
	return s
}
