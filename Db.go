package main

import (
	"path"

	godbhelper "github.com/JojiiOfficial/GoDBHelper"
)

func getDefaultDBFile() string {
	return path.Join(getDataPath(), DefaultDatabaseFile)
}

func connectDB() error {
	dab, err := godbhelper.NewDBHelper(godbhelper.Sqlite).Open(database)
	if err != nil {
		return err
	}
	db = dab
	db.Options.Debug = *appDebug
	db.Options.UseColors = !(*appNoColor)
	return updateDB()
}

func updateDB() error {
	db.AddQueryChain(getInitSQL())
	return db.RunUpdate()
}
