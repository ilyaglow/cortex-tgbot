package cortexbot

import (
	"github.com/boltdb/bolt"
)

func (c *Client) registerUser(u string, method string) {
	c.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(c.UsersBucket))
		err := b.Put([]byte(u), []byte("password"))
		return err
	})
}

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
