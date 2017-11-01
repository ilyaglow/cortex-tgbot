package cortexbot

import (
	"errors"
	"log"
	"math/rand"
	"testing"

	"github.com/boltdb/bolt"
)

var nums = []rune("0123456789")

func generateFilename() string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = nums[rand.Intn(len(nums))]
	}
	return "test" + string(b) + ".db"
}

func mockupClient() *Client {
	db, err := bolt.Open(generateFilename(), 0644, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a bucket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("test"))
		if err != nil {
			return errors.New("Create users bucket failed")
		}
		return nil
	})

	return &Client{
		DB:          db,
		UsersBucket: "test",
	}
}

func TestUser(t *testing.T) {
	c := mockupClient()
	c.registerUser("sample1", "password")
	c.registerUser("aduser", "active_directory")
	c.registerUser("googleuser", "oauth")

	if len(c.listUsers()) != 3 {
		t.Error("There are not 3 users in a bucket as supposed to be")
	}

	if !c.userExists("aduser") {
		t.Error("Registered user doesn't exist")
	}

	if c.userExists("nonexistent") {
		t.Error("Non-existent user exists")
	}

}
