package cortexbot

import (
	"os"
	"strconv"
	"testing"

	tb "gopkg.in/tucnak/telebot.v2"
)

func TestCheckAuth(t *testing.T) {
	c := mockupClient()
	defer os.Remove("test.db")

	var authtest = []struct {
		u           *tb.User
		checkString string
	}{
		{
			&tb.User{
				ID:        1000,
				FirstName: "Name",
				Username:  "user1",
			}, strconv.Itoa(1000)},
		{
			&tb.User{
				ID:        2000,
				FirstName: "Name2",
			}, strconv.Itoa(2000)},
	}

	for _, at := range authtest {
		if err := c.registerUser(at.u); err != nil {
			t.Error(err)
		}

		if !c.CheckAuth(at.u.ID) {
			t.Errorf("check auth for %v failed", at.u)
		}
	}

	if c.CheckAuth(9999) {
		t.Error("Non-existent user bypassed auth check")
	}
}
