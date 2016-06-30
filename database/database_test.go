package database

import (
	"log"
	"os"
	"testing"
)

const dbPath = "/tmp/testing-show.db"

var testDB Database

func removeDatabase() error {
	return os.Remove(dbPath)
}

func setupDatabase() {
	var err error

	removeDatabase() // Make sure there is no already existing database

	testDB, err = Open(dbPath)
	if err != nil {
		log.Fatalf("Open database should be successfull, returned error: %v", err)
	}
}

func TestMain(m *testing.M) {
	populateMockUsers()

	runResult := m.Run()

	removeDatabase() // Clean up the database
	os.Exit(runResult)
}
