package cortexbot

import (
	"errors"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func mockupClient() *Client {
	db, err := bolt.Open("test.db", 0644, nil)
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
	defer os.Remove("test.db")

	var usertests = []struct {
		u  *tgbotapi.User
		id string
	}{
		{&tgbotapi.User{
			ID: 10000,
		}, strconv.Itoa(10000)},
		{&tgbotapi.User{
			ID: 20000,
		}, strconv.Itoa(20000)},
	}

	for _, ut := range usertests {
		if err := c.registerUser(ut.u); err != nil {
			t.Error(err)
		}

		if !c.userExists(ut.id) {
			t.Errorf("%s is not found", ut.id)
		}
	}

	if len(c.listUsers()) != 2 {
		t.Error("There are not 2 users in a bucket as supposed to be")
	}

	if c.userExists("nonexistent") {
		t.Error("Non-existent user exists")
	}
}
