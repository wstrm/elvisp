package database

import (
	"log"
	"reflect"
	"testing"

	"github.com/willeponken/go-cjdns/key"
)

type mockUser struct {
	pubkey  *key.Public
	id      uint64
	invalid bool
}

var binUint64Tests = []struct {
	v uint64
	b []byte
}{
	{1, []byte{0, 0, 0, 0, 0, 0, 0, 1}},
	{18446744073709551616 - 1, []byte{255, 255, 255, 255, 255, 255, 255, 255}},
}

func generateDuplicateUsers() (mockUsers []mockUser) {
	mockKey, err := key.DecodePublic("lpu15wrt3tb6d8vngq9yh3lr4gmnkuv0rgcd2jwl5rp5v0mhlg30.k")

	if err != nil {
		log.Fatalf("populateMockUsers() returned unexpected error: %v", err)
	}

	mockUsers = append(mockUsers, mockUser{mockKey, 1, false}, mockUser{mockKey, 1, true}) // Second is invalid because it's duplicate

	return
}

func generateMockUsers(numUsers int) []mockUser {

	mockUsers := make([]mockUser, numUsers)
	for i := 0; i < numUsers; i++ {
		mockPubKey := key.Generate().Pubkey()

		mockUsers[i] = mockUser{mockPubKey, uint64(i + 1), false}
	}

	return mockUsers
}

func TestUint64ToBin(t *testing.T) {
	for row, test := range binUint64Tests {
		b := uint64ToBin(test.v)

		if !reflect.DeepEqual(b, test.b) {
			t.Errorf("Row: %d returned unexpected binary, got: %v, wanted: %v", row, b, test.b)
		}
	}
}

func TestBinToUint64(t *testing.T) {
	for row, test := range binUint64Tests {
		v := binToUint64(test.b)

		if v != test.v {
			t.Errorf("Row: %d returned unexpected number, got: %d, wanted: %d", row, v, test.v)
		}
	}
}

// TestAddUser_duplicate checks if adding duplicate users returns an error
func TestAddUser_duplicate(t *testing.T) {
	setupDatabase()
	mockUsers := generateDuplicateUsers()

	for row, test := range mockUsers {
		id, err := testDB.AddUser(test.pubkey)

		// Got error, but it should be a valid call
		if err != nil && !test.invalid {
			t.Errorf("Row: %d returned unexpected error: %v", row, err)
		}

		// Got no error, but it shouldn't be a valid call
		if err == nil && test.invalid {
			t.Errorf("Row: %d expected error but got id: %v", row, id)
		}

		// Got wrong id, and the call should be valid (i.e. we expect the correct id)
		if id != test.id && !test.invalid {
			t.Errorf("Row: %d unexpected id, got id: %d, wanted id: %d", row, id, test.id)
		}
	}
}

// TestAddUser_many checks if we're able to add 1000 users without hickups to the database
func TestAddUser_many(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	setupDatabase()
	mockUsers := generateMockUsers(1000)

	for row, test := range mockUsers {
		id, err := testDB.AddUser(test.pubkey)

		// Got error, but it should be a valid call
		if err != nil && !test.invalid {
			t.Errorf("Row: %d returned unexpected error: %v", row, err)
		}

		// Got no error, but it shouldn't be a valid call
		if err == nil && test.invalid {
			t.Errorf("Row: %d expected error but got id: %v", row, id)
		}

		// Got wrong id, and the call should be valid (i.e. we expect the correct id)
		if id != test.id && !test.invalid {
			t.Errorf("Row: %d unexpected id, got id: %d, wanted id: %d", row, id, test.id)
		}
	}
}

// TestAddUser_fillGap checks if the gap created by adding 3 users then removing the 2nd will be filled when adding an user again
func TestAddUser_fillGap(t *testing.T) {
	setupDatabase()
	mockUsers := generateMockUsers(3)

	// Populate the database with the 3 users
	for _, test := range mockUsers {
		testDB.AddUser(test.pubkey)
	}

	secondUser := mockUsers[1]
	testDB.DelUser(secondUser.pubkey)          // Remove the second user
	id, _ := testDB.AddUser(secondUser.pubkey) // Add the user again, should get the same as last time (2)

	if id != secondUser.id {
		t.Errorf("User was not added inside the gap, got id: %d, wanted: %d", id, secondUser.id)
	}
}
