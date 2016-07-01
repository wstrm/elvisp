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
	var err error
	var key1, key2 *key.Public

	key1, err = key.DecodePublic("lpu15wrt3tb6d8vngq9yh3lr4gmnkuv0rgcd2jwl5rp5v0mhlg30.k")
	key2 = key1 // Use same user to test for duplicates (this is why it's "invalid")

	if err != nil {
		log.Fatalf("populateMockUsers() returned unexpected error: %v", err)
	}

	user1 := mockUser{key1, 1, false}
	user2 := mockUser{key2, 1, true}

	mockUsers = append(mockUsers, user1, user2)

	return
}

func generateManyUsers(numUsers int) []mockUser {

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

func TestAddUserDuplicate(t *testing.T) {
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

func TestAddUserMany(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	setupDatabase()
	mockUsers := generateManyUsers(1000)

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
