package cortexbot

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/proxy"

	"github.com/boltdb/bolt"
	"github.com/ilyaglow/go-cortex"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	// defaultTLP is Green because indicators reaching telegram servers
	// TODO: think about making it configurable
	defaultTLP      = 1
	boltFileName    = "bolt.db"
	bucket          = "users"
	tgTokenEnv      = "TGBOT_API_TOKEN"
	cortexURLEnv    = "CORTEX_URL"
	cortexAPIKeyEnv = "CORTEX_API_KEY"
	socksURL        = "SOCKS5_URL"
	cortexBotPWEnv  = "CORTEX_BOT_PASSWORD"
)

var (
	pollTimeout   = 20 * time.Second
	cortexTimeout = 5 * time.Minute
)

// Client defines bot's abilities to interact with services
type Client struct {
	Bot         *tb.Bot
	Cortex      *cortex.Client
	Password    string
	DB          *bolt.DB
	UsersBucket string
	TLP         int
	Timeout     time.Duration
	BotSettings *tb.Settings
}

// socks5Client bootstraps http.Client that uses socks5 proxy
func socks5Client(u *url.URL) (*http.Client, error) {
	dialer, err := proxy.FromURL(u, proxy.Direct)
	if err != nil {
		return nil, err
	}

	return &http.Client{Transport: &http.Transport{Dial: dialer.Dial}}, nil
}

func usage() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	fmt.Printf(`Usage:
	%s=<telegram bot token> \
	%s=<cortex location> \
	%s=<cortex API key> \
	%s=<cortex bot passphrase> \
	%s
`, tgTokenEnv, cortexURLEnv, cortexAPIKeyEnv, cortexBotPWEnv, ex)
	os.Exit(1)
}

// NewClient bootstraps the Client struct from env variables
func NewClient() *Client {
	var (
		bot *tb.Bot
		err error
	)

	if os.Getenv(tgTokenEnv) == "" || os.Getenv(cortexURLEnv) == "" || os.Getenv(cortexAPIKeyEnv) == "" || os.Getenv(cortexBotPWEnv) == "" {
		usage()
	}

	tgToken := os.Getenv(tgTokenEnv)

	var settings tb.Settings
	if proxy, ok := os.LookupEnv("SOCKS5_URL"); ok {
		surl, err := url.Parse(proxy)
		if err != nil {
			log.Panic(err)
		}

		sc, err := socks5Client(surl)
		if err != nil {
			log.Panic(err)
		}

		settings = tb.Settings{
			Token: tgToken,
			Poller: &tb.LongPoller{
				Timeout: pollTimeout,
			},
			Client: sc,
		}
	} else {
		settings = tb.Settings{
			Token: tgToken,
			Poller: &tb.LongPoller{
				Timeout: pollTimeout,
			},
		}
	}

	bot, err = tb.NewBot(settings)
	if err != nil {
		log.Fatal(err)
	}

	crtx, err := cortex.NewClient(os.Getenv(cortexURLEnv), &cortex.ClientOpts{
		Auth: &cortex.APIAuth{
			APIKey: os.Getenv(cortexAPIKeyEnv),
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
		timeout = cortexTimeout
	} else {
		timeout, errt = time.ParseDuration(timeoutStr)
		if errt != nil {
			log.Fatal(errt)
		}
	}

	return &Client{
		Bot:         bot,
		Cortex:      crtx,
		Password:    os.Getenv(cortexBotPWEnv),
		DB:          db,
		UsersBucket: bucket,
		TLP:         defaultTLP,
		Timeout:     timeout,
		BotSettings: &settings,
	}
}
