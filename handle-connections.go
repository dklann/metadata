// handle-connections.go contains connection handlers.
package main

import (
	"encoding/xml"
	"io"
	"log"
	"net"

	"fmt"
)

type Metadata struct {
	Time   string `xml:"time"`
	Artist string `xml:"artist"`
	Title  string `xml:"title"`
	Album  string `xml:"album"`
}

// String converts the Metadata struct into a printable string.
func (m Metadata) String() string {
	s := fmt.Sprintf("m: time: %s artist: %s, title: %s, album: %s",
		m.Time,
		m.Artist,
		m.Title,
		m.Album)
	return s
}

// handleConnection snags the XML text from the remote,
// unmarshals it and hands it off to the next step...
func handleConnection(connection net.Conn) error {
	var metadata Metadata

	defer connection.Close()
	data, err := io.ReadAll(connection)
	if err != nil {
		log.Println(logSprintf("unable to ReadAll from '%v' (%v).", connection, err))
		return err
	}

	err = xml.Unmarshal(data, &metadata)
	if err != nil {
		log.Println(logSprintf("error unmarshaling %v (%v).", data, err))
		return err
	}
	verbosePrint(fmt.Sprintf("metadata: %s\n", metadata.String()))

	db, err := dbOpen("m")
	if err != nil {
		log.Println(logSprintf("error opening the database; better luck next time: %v", err))
		return err
	}
	defer db.Close()

	err = metadata.writeRecord(db)
	if err != nil {
		log.Println(logSprintf("error writing metadata to the database: %v", err))
		return err
	}

	err = ageDatabase(db)
	if err != nil {
		log.Println(logSprintf("error deleting old records from the database: %v", err))
		return err
	}

	return nil
}
