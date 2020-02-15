package main

import (
	"database/sql"
	"fmt"
	"strconv"

	dbhelper "github.com/JojiiOfficial/GoDBHelper"
)

//Tables
const (
	//TableSubscriptions table for subscriptions
	TableSubscriptions = "subscriptions"
	//TableActions table for actions
	TableActions = "actions"
)

func getInitSQL() dbhelper.QueryChain {
	return dbhelper.QueryChain{
		Name:  "initChain",
		Order: 0,
		Queries: dbhelper.CreateInitVersionSQL(
			dbhelper.InitSQL{
				Query:   "CREATE TABLE %s (`pkID` INTEGER PRIMARY KEY AUTOINCREMENT, `sourceID` TEXT, `hookName` TEXT, `subsID` TEXT, `valid` INTEGER DEFAULT 0)",
				FParams: []string{TableSubscriptions},
			},
			dbhelper.InitSQL{
				Query:   "CREATE TABLE %s (`pkID` INTEGER PRIMARY KEY AUTOINCREMENT, `name` TEXT, `subscriptionID` INTEGER, `mode` INTEGER, `file` TEXT)",
				FParams: []string{TableActions},
			},
		),
	}
}

//Global sql queries

// ---------------------------------- Inserts ------------------------------------

func (webhook *Subscription) insert(db *dbhelper.DBhelper) error {
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

func (action *Action) insert(db *dbhelper.DBhelper) error {
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

//-->  Actions ----------------
func getActionIDs(db *dbhelper.DBhelper, limit int) ([]string, error) {
	var hooks []string
	err := db.QueryRowsf(&hooks, "SELECT pkID FROM %s LIMIT %s", []string{TableActions, strconv.Itoa(limit)})
	return hooks, err
}

func hasAction(db *dbhelper.DBhelper, actionName string) (bool, error) {
	var c int
	err := db.QueryRowf(&c, "SELECT COUNT(*) FROM %s WHERE name=?", []string{TableActions}, actionName)
	if err != nil {
		return false, err
	}
	return c > 0, nil
}

func getActionIDFromName(db *dbhelper.DBhelper, actionName string) (int64, error) {
	var c int64
	err := db.QueryRowf(&c, "SELECT pkID FROM %s WHERE name=?", []string{TableActions}, actionName)
	if err != nil {
		return -1, err
	}
	return c, nil
}

func getActions(db *dbhelper.DBhelper) ([]Action, error) {
	var actions []Action
	err := db.QueryRowsf(&actions, "SELECT %s.*,IFNULL(%s.hookName,'- na -')as hookName FROM %s LEFT JOIN %s ON (%s.pkID = %s.subscriptionID) ORDER BY pkID DESC", []string{TableActions, TableSubscriptions, TableActions, TableSubscriptions, TableSubscriptions, TableActions})
	return actions, err
}

func getActionFromName(db *dbhelper.DBhelper, actionName string) (*Action, error) {
	var action Action
	err := db.QueryRowf(&action, "SELECT * FROM %s WHERE name=?", []string{TableActions}, actionName)
	if err != nil {
		return nil, err
	}
	return &action, nil
}

func getActionsForSource(db *dbhelper.DBhelper, sourceID int64) ([]Action, error) {
	var actions []Action
	err := db.QueryRowsf(&actions, "SELECT * FROM %s WHERE subscriptionID=? ORDER BY pkID ASC", []string{TableActions}, sourceID)
	return actions, err
}

//-->  Subscriptions ---------------------
func getSubscriptions(db *dbhelper.DBhelper) ([]Subscription, error) {
	var subscriptions []Subscription
	err := db.QueryRowsf(&subscriptions, "SELECT * FROM %s ORDER BY pkID DESC", []string{TableSubscriptions})
	return subscriptions, err
}

func getSubscriptionsHumanized(db *dbhelper.DBhelper, limit int) ([]string, error) {
	var hooks []string
	err := db.QueryRowsf(&hooks, "SELECT hookName || '-' || sourceID FROM %s WHERE hookName != ''", []string{TableSubscriptions})
	return hooks, err
}

func getSubscriptionFromID(db *dbhelper.DBhelper, sid int64) (*Subscription, error) {
	var subscription Subscription
	err := db.QueryRowf(&subscription, "SELECT * FROM %s WHERE pkID=?", []string{TableSubscriptions}, sid)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func getSubscriptionFromSourceID(db *dbhelper.DBhelper, sid string) (*Subscription, error) {
	var subscription Subscription
	err := db.QueryRowf(&subscription, "SELECT * FROM %s WHERE sourceID=?", []string{TableSubscriptions}, sid)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func getSubscriptionID(db *dbhelper.DBhelper, sourceID string) (int64, error) {
	var id int64
	err := db.QueryRowf(&id, "SELECT pkID FROM %s WHERE sourceID=?", []string{TableSubscriptions}, sourceID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func deleteSubscription(db *dbhelper.DBhelper, pkID int64) error {
	_, err := db.Execf("DELETE FROM %s WHERE pkID=?", []string{TableSubscriptions}, pkID)
	return err
}

func deleteSubscriptionByID(db *dbhelper.DBhelper, sourceID string) error {
	_, err := db.Execf("DELETE FROM %s WHERE sourceID=?", []string{TableSubscriptions}, sourceID)
	return err
}

func hasSubscribted(db *dbhelper.DBhelper, subscriptionID, sourceID string) (bool, error) {
	var c int
	err := db.QueryRowf(&c, "SELECT COUNT(*) FROM %s WHERE subsID=? AND sourceID=?", []string{TableSubscriptions}, subscriptionID, sourceID)
	if err != nil {
		return false, err
	}
	return c > 0, nil
}

func validateSubscription(db *dbhelper.DBhelper, subscriptionID string) error {
	_, err := db.Execf("UPDATE %s SET valid=1 WHERE subsID=?", []string{TableSubscriptions}, subscriptionID)
	fmt.Printf("UPDATE %s SET valid=1 WHERE subsID=%s", TableSubscriptions, subscriptionID)
	return err
}

// ---------------------------------- Updates ------------------------------------

func removeActionSource(db *dbhelper.DBhelper, subscriptionID int64) error {
	_, err := db.Execf("UPDATE %s SET subscriptionID=0 WHERE subscriptionID=?", []string{TableActions}, subscriptionID)
	return err
}
func updateActionSource(db *dbhelper.DBhelper, aID, subscriptionID int64) error {
	_, err := db.Execf("UPDATE %s SET subscriptionID=? WHERE pkID=?", []string{TableActions}, subscriptionID, aID)
	return err
}

func updateActionMode(db *dbhelper.DBhelper, aID int64, newMode uint8) error {
	_, err := db.Execf("UPDATE %s SET mode=? WHERE pkID=?", []string{TableActions}, newMode, aID)
	return err
}

func updateActionFile(db *dbhelper.DBhelper, aID int64, newFile string) error {
	_, err := db.Execf("UPDATE %s SET file=? WHERE pkID=?", []string{TableActions}, newFile, aID)
	return err
}

// ---------------------------------- Deletions ------------------------------------

func deleteActionByID(db *dbhelper.DBhelper, name string) error {
	_, err := db.Execf("DELETE FROM %s WHERE name=?", []string{TableActions}, name)
	return err
}
