package main

import (
	"database/sql"
	"strconv"

	godbhelper "github.com/JojiiOfficial/GoDBHelper"
)

//Tables
const (
	//TableSubscriptions table for subscriptions
	TableSubscriptions = "subscriptions"
	//TableActions table for actions
	TableActions = "actions"
)

func getInitSQL() godbhelper.QueryChain {
	return godbhelper.QueryChain{
		Name:  "initChain",
		Order: 0,
		Queries: godbhelper.CreateInitVersionSQL(
			godbhelper.InitSQL{
				Query:   "CREATE TABLE %s (`pkID` INTEGER PRIMARY KEY AUTOINCREMENT, `hookID` TEXT, `hookName` TEXT, `subsID` TEXT)",
				FParams: []string{TableSubscriptions},
			},
			godbhelper.InitSQL{
				Query:   "CREATE TABLE %s (`pkID` INTEGER PRIMARY KEY AUTOINCREMENT, `name` TEXT, `subscriptionID` INTEGER, `mode` INTEGER, `file` TEXT)",
				FParams: []string{TableActions},
			},
		),
	}
}

//Global sql queries

// ---------------------------------- Inserts ------------------------------------

func (webhook *Subscription) insert(db *godbhelper.DBhelper) error {
	var rs *sql.Result
	var err error
	var id int64

	if rs, err = db.Insert(*webhook, TableSubscriptions); err != nil || rs == nil {
		return err
	}

	if id, err = (*rs).LastInsertId(); err != nil {
		return err
	}

	webhook.ID = id
	return nil
}

func (action *Action) insert(db *godbhelper.DBhelper) error {
	var rs *sql.Result
	var err error
	var id int64

	if rs, err = db.Insert(*action, TableActions); err != nil || rs == nil {
		return err
	}

	if id, err = (*rs).LastInsertId(); err != nil {
		return err
	}

	action.ID = id
	return nil
}

// ---------------------------------- Selects ------------------------------------

func getActions(db *godbhelper.DBhelper) ([]Action, error) {
	var actions []Action
	err := db.QueryRowsf(&actions, "SELECT %s.*,IFNULL(%s.hookName,'- na -')as hookName FROM %s LEFT JOIN %s ON (%s.pkID = %s.subscriptionID) ORDER BY pkID DESC", []string{TableActions, TableSubscriptions, TableActions, TableSubscriptions, TableSubscriptions, TableActions})
	return actions, err
}

func getSubscriptions(db *godbhelper.DBhelper) ([]Subscription, error) {
	var subscriptions []Subscription
	err := db.QueryRowsf(&subscriptions, "SELECT * FROM %s ORDER BY pkID DESC", []string{TableSubscriptions})
	return subscriptions, err
}

func getHooksHumanized(db *godbhelper.DBhelper, limit int) ([]string, error) {
	var hooks []string
	err := db.QueryRowsf(&hooks, "SELECT hookName || '-' || hookID FROM %s WHERE hookName != ''", []string{TableSubscriptions})
	return hooks, err
}

func getSubscriptionID(db *godbhelper.DBhelper, sid string) (int64, error) {
	var id int64
	err := db.QueryRowf(&id, "SELECT pkID FROM %s WHERE hookID=?", []string{TableSubscriptions}, sid)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func getActionIDs(db *godbhelper.DBhelper, limit int) ([]string, error) {
	var hooks []string
	err := db.QueryRowsf(&hooks, "SELECT pkID FROM %s LIMIT %s", []string{TableActions, strconv.Itoa(limit)})
	return hooks, err
}

func hasAction(db *godbhelper.DBhelper, actionName string) (bool, error) {
	var c int
	err := db.QueryRowf(&c, "SELECT COUNT(*) FROM %s WHERE name=?", []string{TableActions}, actionName)
	if err != nil {
		return false, err
	}
	return c > 0, nil
}

func getActionFromName(db *godbhelper.DBhelper, actionName string) (int64, error) {
	var c int64
	err := db.QueryRowf(&c, "SELECT pkID FROM %s WHERE name=?", []string{TableActions}, actionName)
	if err != nil {
		return -1, err
	}
	return c, nil
}
func updateActionWebhook(db *godbhelper.DBhelper, aID, whID int64) error {
	_, err := db.Execf("UPDATE %s SET subscriptionID=? WHERE pkID=?", []string{TableActions}, whID, aID)
	return err
}

// ---------------------------------- Deletions ------------------------------------

func deleteActionByID(db *godbhelper.DBhelper, name string) error {
	_, err := db.Execf("DELETE FROM %s WHERE name=?", []string{TableActions}, name)
	return err
}
