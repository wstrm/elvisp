package database_test

import (
	"io/ioutil"
	"os"

	"github.com/willeponken/elvisp/database"
)

func tempFile() string {
	file, err := ioutil.TempFile("", "elvisp-")
	if err != nil {
		panic(err)
	}

	if err := file.Close(); err != nil {
		panic(err)
	}

	if err := os.Remove(file.Name()); err != nil {
		panic(err)
	}

	return file.Name()
}

type TestDB struct {
	database.Database
}

func MustOpen() TestDB {
	db, err := database.Open(tempFile())
	if err != nil {
		panic(err)
	}

	return TestDB{db}
}

func (t *TestDB) Close() error {
	defer os.Remove(t.Path())
	return t.Database.Close()
}

func (t *TestDB) MustClose() {
	if err := t.Close(); err != nil {
		panic(err)
	}
}
