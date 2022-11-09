// process-data.go processes the metadata retrieved in handleConnection().
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/araddon/dateparse"
)

// getTime gets the time into a usable format from the data.
func (m Metadata) getTime() time.Time {
	dateTime, err := dateparse.ParseAny(m.Time)
	if err != nil {
		log.Println(logSprintf("unable to parse '%v% (%v)", m.Time, err))
		return time.Time{}
	}
	verbosePrint(fmt.Sprintf("metadata timestamp: %v", dateTime))
	return dateTime
}

// writeRecord writes the supplied metadata to the database.
func (m Metadata) writeRecord(db *sql.DB) error {
	q := "INSERT INTO `m` "
	q += "(time, artist, title, album) "
	q += "VALUES "
	q += "(?,    ?,      ?,     ?)"
	verbosePrint(fmt.Sprintf("about to Prepare statement: %s\n", q))
	stmt, err := db.Prepare(q)
	if err != nil {
		log.Println(logSprintf("unable to prepare statement: %#v.", err))
		return err
	}
	defer stmt.Close()
	verbosePrint(fmt.Sprintf("inserting fields: %s\n", m.String()))
	var result sql.Result
	result, err = stmt.Exec(
		m.getTime().UTC().Format("2006-01-02 15:04:05"),
		m.Artist,
		m.Title,
		m.Album)
	if err != nil {
		log.Println(logSprintf("error in stmt.Exec(): %#v", err))
		return err
	}
	debugPrint(fmt.Sprintf("INSERT result: %v\n", result))

	return nil
}

// ageDatabase deletes records from the database
// that are older than 1 year old.
func ageDatabase(db *sql.DB) error {
	q := "DELETE FROM `m` WHERE time < NOW() - INTERVAL 1 YEAR"
	verbosePrint(fmt.Sprintf("about to Prepare statement: %s\n", q))
	stmt, err := db.Prepare(q)
	if err != nil {
		log.Println(logSprintf("unable to prepare statement: %#v.", err))
		return err
	}
	defer stmt.Close()
	var result sql.Result
	result, err = stmt.Exec()
	if err != nil {
		log.Println(logSprintf("error in stmt.Exec(): %#v", err))
		return err
	}
	debugPrint(fmt.Sprintf("INSERT result: %v\n", result))
	return nil
}
