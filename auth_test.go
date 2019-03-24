package cortexbot

import (
	"os"
	"testing"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestCheckAuth(t *testing.T) {
	c := mockupClient()
	defer os.Remove("test.db")

	var authtest = []struct {
		tgUser *tgbotapi.User
		user   *User
	}{
		{
			&tgbotapi.User{
				ID:        1000,
				FirstName: "Name",
				UserName:  "user1",
			},
			&User{
				ID: 1000,
			},
		},
		{
			&tgbotapi.User{
				ID:        2000,
				FirstName: "Name2",
			},
			&User{
				ID: 2000,
			},
		},
	}

	for _, at := range authtest {
		if err := c.addUser(at.user); err != nil {
			t.Error(err)
		}

		if !c.CheckAuth(at.tgUser) {
			t.Errorf("check auth for %v failed", at.tgUser)
		}
	}

	if c.CheckAuth(&tgbotapi.User{
		ID: 4000}) {
		t.Error("Non-existent user bypassed auth check")
	}
}
