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
		s = append(s, e.Name())
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
	database = getDBunparsed()
	if err := connectDB(); err != nil {
		fmt.Println(err.Error())
	}
	dat, err := getActionIDs(db, 9)
	if err != nil {
		log.Fatalln(err.Error())
		return []string{}
	}
	return dat
}

func hintSubscriptions() []string {
	database = getDBunparsed()
	if err := connectDB(); err != nil {
		fmt.Println(err.Error())
	}
	dat, err := getHooksHumanized(db, 9)
	if err != nil {
		log.Fatalln(err.Error())
		return []string{}
	}
	return dat
}
