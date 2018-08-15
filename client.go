package cortexbot

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/proxy"

	"github.com/boltdb/bolt"
	"github.com/ilyaglow/go-cortex"
	"github.com/ilyaglow/telegram-bot-api"
)

const (
	// defaultTLP is Green because indicators reaching telegram servers
	// TODO: think about making it configurable
	defaultTLP   = 1
	boltFileName = "bolt.db"
	bucket       = "users"
)

var defaultTimeout = 5 * time.Minute

// Client defines bot's abilities to interact with services
type Client struct {
	Bot         *tgbotapi.BotAPI
	Cortex      *cortex.Client
	Password    string
	DB          *bolt.DB
	UsersBucket string
	TLP         int
	Timeout     time.Duration
}

// socks5Client bootstraps http.Client that uses socks5 proxy
func socks5Client(u *url.URL) (*http.Client, error) {
	dialer, err := proxy.FromURL(u, proxy.Direct)
	if err != nil {
		return nil, err
	}

	return &http.Client{Transport: &http.Transport{Dial: dialer.Dial}}, nil
}

// NewClient bootstraps the Client struct from env variables
func NewClient() *Client {
	var (
		bot *tgbotapi.BotAPI
		err error
	)

	tgToken, ok := os.LookupEnv("TGBOT_API_TOKEN")
	if !ok {
		log.Fatal("TGBOT_API_TOKEN environment variable is not set")
	}

	if proxy, ok := os.LookupEnv("SOCKS5_URL"); ok {
		surl, err := url.Parse(proxy)
		if err != nil {
			log.Panic(err)
		}

		sc, err := socks5Client(surl)
		if err != nil {
			log.Panic(err)
		}

		bot, err = tgbotapi.NewBotAPIWithClient(tgToken, sc)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		bot, err = tgbotapi.NewBotAPI(tgToken)
		if err != nil {
			log.Fatal(err)
		}
	}

	crtx, err := cortex.NewClient(os.Getenv("CORTEX_LOCATION"), &cortex.ClientOpts{
		Auth: &cortex.APIAuth{
			APIKey: os.Getenv("CORTEX_API_KEY"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	db, err := bolt.Open("bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a bucket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return errors.New("Create users bucket failed")
		}
		return nil
	})

	var (
		timeout time.Duration
		errt    error
	)
	timeoutStr, ok := os.LookupEnv("CORTEX_TIMEOUT")
	if !ok {
		timeout = defaultTimeout
	} else {
		timeout, errt = time.ParseDuration(timeoutStr)
		if errt != nil {
			log.Fatal(errt)
		}
	}

	return &Client{
		Bot:         bot,
		Cortex:      crtx,
		Password:    os.Getenv("CORTEX_BOT_PASSWORD"),
		DB:          db,
		UsersBucket: bucket,
		TLP:         defaultTLP,
		Timeout:     timeout,
	}
}
