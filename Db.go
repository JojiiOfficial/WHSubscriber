package main

import (
	"path"

	dbhelper "github.com/JojiiOfficial/GoDBHelper"
)

func getDefaultDBFile() string {
	return path.Join(getDataPath(), DefaultDatabaseFile)
}

func connectDB(dbFile string) (*dbhelper.DBhelper, error) {
	db, err := dbhelper.NewDBHelper(dbhelper.Sqlite).Open(dbFile)
	if err != nil {
		return nil, err
	}
	db.Options.Debug = *appDebug
	db.Options.UseColors = !(*appNoColor)
	return db, updateDB(db)
}

func updateDB(db *dbhelper.DBhelper) error {
	db.AddQueryChain(getInitSQL())
	return db.RunUpdate()
}
