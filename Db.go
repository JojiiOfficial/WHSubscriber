package main

import (
	"path"

	godbhelper "github.com/JojiiOfficial/GoDBHelper"
)

func getDefaultDBFile() string {
	return path.Join(getDataPath(), DefaultDatabaseFile)
}

func connectDB(dbFile string) (*godbhelper.DBhelper, error) {
	db, err := godbhelper.NewDBHelper(godbhelper.Sqlite).Open(dbFile)
	if err != nil {
		return nil, err
	}
	db.Options.Debug = *appDebug
	db.Options.UseColors = !(*appNoColor)
	return db, updateDB(db)
}

func updateDB(db *godbhelper.DBhelper) error {
	db.AddQueryChain(getInitSQL())
	return db.RunUpdate()
}
