package cortexbot

import (
	"database/sql"
	"log"
)

// // registerUser adds a user to boltdb bucket.
// func (c *Cortexbot) registerUser(u *tgbotapi.User) error {
// 	data, err := json.Marshal(u)
// 	if err != nil {
// 		return err
// 	}

// 	c.DB.Update(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte(c.UsersBucket))
// 		err = b.Put([]byte(strconv.Itoa(u.ID)), data)
// 		return err
// 	})

// 	return err
// }

// // listUsers returns a slice of all users
// // Users means keys in a users bucket
// func (c *Cortexbot) listUsers() []string {
// 	var users []string

// 	c.DB.View(func(tx *bolt.Tx) error {
// 		// Assume bucket exists and has keys
// 		b := tx.Bucket([]byte(c.UsersBucket))

// 		c := b.Cursor()

// 		for k, _ := c.First(); k != nil; k, _ = c.Next() {
// 			users = append(users, string(k))
// 		}

// 		return nil
// 	})

// 	return users
// }

// // userExists checks if user exists among registered users
// func (c *Cortexbot) userExists(u string) bool {
// 	exists := false

// 	c.DB.View(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte(c.UsersBucket))
// 		v := b.Get([]byte(u))
// 		if v != nil {
// 			exists = true
// 		}
// 		return nil
// 	})

// 	return exists
// }

// func (c *Cortexbot) addAdminID(chatID int) error {
// 	c.DB.Update(func(tx *bolt.Tx) error {
// 		err = tx.Bucket([]byte(c.OptsBucket)).Put([]byte("admin_chat_id")).strconv.Itoa(chatID)
// 		return err
// 	}
// 	return err
// }

// func (c *Cortexbot) getChatWithAdmin() (chatID int64, err error) {
// 	err = c.DB.View(func(tx *bolt.Tx) error {
// 		v := tx.Bucket([]byte(c.OptsBucket)).Get([]byte("admin_chat_id"))
// 		if v == nil {
// 			return errors.New("no admin registered")
// 		}

// 		buf := bytes.NewBuffer(v)
// 		binary.Read(buf, binary.LittleEndian, &chatID)

// 		return nil
// 	})
// 	return
// }

func (c *Cortexbot) createUsersTbl() error {
	stmt, err := c.DB.Prepare(`
		CREATE TABLE IF NOT EXISTS users(id INTEGER PRIMARY KEY, admin INTEGER, info TEXT)
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

// User represents a user of the bot.
type User struct {
	ID    int
	Admin int // sqlite doesn't have bool
	About string
}

func (c *Cortexbot) getUser(id int) (*User, error) {
	var user User
	err := c.DB.QueryRow(`
		SELECT id, admin, info
		FROM users
		WHERE id=?
	`, id).Scan(&user.ID, &user.Admin, &user.About)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Cortexbot) getAdmins() ([]*User, error) {
	var admins []*User
	rows, err := c.DB.Query(`
		SELECT id, info
		FROM users
		WHERE admin=?
	`, 1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		u := User{
			Admin: 1,
		}
		err := rows.Scan(&u.ID, &u.About)
		if err != nil {
			return nil, err
		}
		admins = append(admins, &u)
	}
	if err = rows.Err(); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return admins, nil
}

func (c *Cortexbot) addUser(u *User) error {
	stmt, err := c.DB.Prepare(`
		INSERT INTO users(id, admin, info) VALUES(?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(u.ID, u.Admin, u.About)
	return err
}

func (c *Cortexbot) delUser(id int) error {
	stmt, err := c.DB.Prepare(`
		DELETE FROM users
		WHERE id=?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	return err
}

func (c *Cortexbot) listUsers() ([]*User, error) {
	var users []*User
	rows, err := c.DB.Query(`
		SELECT id, admin, info
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Admin, &u.About)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	if err = rows.Err(); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return users, nil
}

func (c *Cortexbot) noAdmin() bool {
	admins, err := c.getAdmins()
	if err != nil {
		return false
	}

	if len(admins) == 0 {
		return true
	}
	return false
}

func (c *Cortexbot) userExists(id int) bool {
	row := c.DB.QueryRow(`
		SELECT 1
		FROM users
		WHERE id=?
	`, id)

	var v int
	if err := row.Scan(&v); err != nil {
		log.Println(err)
		return false
	}

	return true
}
