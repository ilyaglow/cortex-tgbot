package cortexbot

import (
	"errors"
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ilyaglow/go-cortex"
)

var boltFileName string = "bolt.db"
var bucket string = "users"

// Client defines bot's abilities to interact with services
type Client struct {
	Bot         *tgbotapi.BotAPI
	Cortex      *gocortex.Client
	Password    string
	DB          *bolt.DB
	UsersBucket string
}

// NewClient bootstraps the Client struct from env variables
func NewClient() *Client {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGBOT_API_TOKEN"))
	if err != nil {
		log.Println("TGBOT_API_TOKEN env variable is empty")
		log.Panic(err)
	}

	cortex := gocortex.NewClient(os.Getenv("CORTEX_LOCATION"))

	db, err := bolt.Open("bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a bucket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(bucket))
		if err != nil {
			return errors.New("Create users bucket failed")
		}
		return nil
	})

	return &Client{
		Bot:         bot,
		Cortex:      cortex,
		Password:    os.Getenv("CORTEX_BOT_PASSWORD"),
		DB:          db,
		UsersBucket: bucket,
	}
}
