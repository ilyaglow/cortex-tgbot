package cortexbot

import (
	"github.com/boltdb/bolt"
)

// registerUser adds a user to boltdb bucket
func (c *Client) registerUser(u string, method string) {
	c.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(c.UsersBucket))
		err := b.Put([]byte(u), []byte("password"))
		return err
	})
}

// listUsers returns a slice of all users
// Users means keys in a users bucket
func (c *Client) listUsers() []string {
	var users []string

	c.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(c.UsersBucket))

		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			users = append(users, string(k))
		}

		return nil
	})

	return users
}

// userExists checks if user exists among registered users
func (c *Client) userExists(u string) bool {
	exists := false

	c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(c.UsersBucket))
		v := b.Get([]byte(u))
		if v != nil {
			exists = true
		}
		return nil
	})

	return exists
}
