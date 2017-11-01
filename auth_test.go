package cortexbot

import (
	"os"
	"testing"
)

func TestCheckAuth(t *testing.T) {
	c := mockupClient()
	defer os.Remove("test.db")

	c.registerUser("user", "password")
	if !c.CheckAuth("user") {
		t.Error("User doesn't exist")
	}

	if c.CheckAuth("nonexistent") {
		t.Error("Non-existent user bypassed auth check")
	}
}
