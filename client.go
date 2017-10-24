package cortexbot

import (
	"log"
	"os"

	"github.com/ilyaglow/go-cortex"
	"gopkg.in/telegram-bot-api.v4"
)

// Client defines bot's abilities to interact with services
type Client struct {
	Bot              *tgbotapi.BotAPI
	Cortex           *gocortex.Client
	Password         string
	AllowedUsernames map[string]bool
}

// NewClient bootstraps the Client struct from env variables
func NewClient() *Client {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGBOT_API_TOKEN"))
	if err != nil {
		log.Println("TGBOT_API_TOKEN env variable is empty")
		log.Panic(err)
	}

	cortex := gocortex.NewClient(os.Getenv("CORTEX_LOCATION"))

	return &Client{
		Bot:              bot,
		Cortex:           cortex,
		Password:         os.Getenv("CORTEX_BOT_PASSWORD"),
		AllowedUsernames: make(map[string]bool),
	}
}
