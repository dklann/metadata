package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/dklann/mycnf"
	_ "github.com/go-sql-driver/mysql"
)

// myCnf is used all over the app for the database info.
var myCnf map[string]string

// dbOpen opens the database and returns a handle and an error.
// Remember to close the database (db.Close()) when we are done with it.
func dbOpen(dbName string) (*sql.DB, error) {
	var myDSN string
	var err error
	var db *sql.DB
	myCnf, err = mycnf.ReadMyCnf(&cmdline.MyConfig, "metadata")
	if err != nil {
		log.Println(logSprintf("error in mycnf.ReadMyCnf(%s, `client`): %#v.", cmdline.MyConfig, err))
	}
	if myCnf == nil {
		myCnf = make(map[string]string, 5)
	} else {
		verbosePrint(fmt.Sprintf("myCnf: %v\n", myCnf))
	}
	// Command line options take precedence over .my.cnf settings.
	if cmdline.DbHost != "" {
		myCnf["host"] = cmdline.DbHost
	}
	if cmdline.DbPort != 0 {
		myCnf["port"] = strconv.Itoa(int(cmdline.DbPort))
	}
	if cmdline.DbName != "" {
		myCnf["database"] = cmdline.DbName
	}
	if cmdline.DbUser != "" {
		myCnf["user"] = cmdline.DbUser
	}
	if cmdline.DbPass != "" {
		myCnf["password"] = cmdline.DbPass
	}
	if myCnf == nil {
		// Use defaults if no .my.cnf config
		myDSN = cmdline.DbUser + ":" + cmdline.DbPass + "@tcp(" + cmdline.DbHost + strconv.Itoa(int(cmdline.DbPort)) + cmdline.DbName
	} else {
		myDSN = myCnf["user"] + ":" + myCnf["password"] + "@tcp(" + myCnf["host"] + ":" + myCnf["port"] + ")/" + myCnf["database"]
	}
	myDSN += "?parseTime=true"
	verbosePrint(fmt.Sprintf("myDSN: %s\n", myDSN))
	db, err = sql.Open("mysql", myDSN)
	if err != nil {
		log.Println(logSprintf("error attempting to open database '%s' (%#v).", dbName, err))
		return nil, err
	}
	// From https://github.com/golang/go/wiki/SQLInterface
	// "Note that Open does not directly open a database connection"
	if err = db.Ping(); err != nil {
		log.Println(logSprintf("error pinging the database (%#v).", err))
		return nil, err
	}
	verbosePrint(fmt.Sprintln("yay, opened and pinged the database"))
	return db, nil
}
