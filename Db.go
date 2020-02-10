package main

import (
	godbhelper "github.com/JojiiOfficial/GoDBHelper"
)

func connectDB() error {
	dab, err := godbhelper.NewDBHelper(godbhelper.Sqlite).Open(database)
	if err != nil {
		return err
	}
	db = dab
	return updateDB()
}

func updateDB() error {
	return db.RunUpdate()
}
