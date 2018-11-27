package cortexbot

import (
	"os"
	"strconv"
	"testing"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestCheckAuth(t *testing.T) {
	c := mockupClient()
	defer os.Remove("test.db")

	var authtest = []struct {
		u           *tgbotapi.User
		checkString string
	}{
		{
			&tgbotapi.User{
				ID:        1000,
				FirstName: "Name",
				UserName:  "user1",
			}, strconv.Itoa(1000)},
		{
			&tgbotapi.User{
				ID:        2000,
				FirstName: "Name2",
			}, strconv.Itoa(2000)},
	}

	for _, at := range authtest {
		if err := c.registerUser(at.u); err != nil {
			t.Error(err)
		}

		if !c.CheckAuth(at.u) {
			t.Errorf("check auth for %v failed", at.u)
		}
	}

	if c.CheckAuth(&tgbotapi.User{
		ID: 4000}) {
		t.Error("Non-existent user bypassed auth check")
	}
}
