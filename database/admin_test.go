package database

import (
	"strconv"
	"testing"
)

// TestSetAdmin_AdminHash_replace checks if adding a hash again replaces the old one, it tests both the SetAdmin and AdminHash methods
func TestSetAdmin_AdminHash_replace(t *testing.T) {
	setupDatabase()

	var currHash string
	var retrHash string
	var err error
	for i := 0; i < 2; i++ {
		currHash = "hashNumber" + strconv.Itoa(i)

		err = testDB.SetAdmin(currHash)
		if err != nil {
			t.Errorf("SetAdmin returned unexpected error: %v", err)
		}

		retrHash, err = testDB.AdminHash()
		if err != nil {
			t.Errorf("SetAdmin returned unexpected error: %v", err)
		}

		if retrHash != currHash {
			t.Errorf("SetAdmin returned unexpected hash, got: %s, wanted: %s", retrHash, currHash)
		}
	}
}
