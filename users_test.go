package cortexbot

import (
	"database/sql"
	"log"
	"os"
	"testing"
)

func mockupClient() *Cortexbot {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}

	c := &Cortexbot{
		DB: db,
	}
	err = c.createUsersTbl()
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func TestUser(t *testing.T) {
	c := mockupClient()
	defer os.Remove("test.db")

	var usertests = []struct {
		u  *User
		id int
	}{
		{&User{
			ID: 10000,
		}, 10000},
		{&User{
			ID: 20000,
		}, 20000},
	}

	for _, ut := range usertests {
		if err := c.addUser(ut.u); err != nil {
			t.Error(err)
		}

		user, err := c.getUser(ut.id)
		if err != nil {
			t.Fatal(err)
		}
		if user == nil {
			t.Errorf("%d is not found", ut.id)
		}
	}

	users, err := c.listUsers()
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Error("there are not 2 users in a bucket as supposed to be")
	}

	nonexistent, err := c.getUser(0)
	if err != nil {
		log.Fatal(err)
	}
	if nonexistent != nil {
		t.Errorf("non-existent account exists")
	}
}
