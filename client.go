package cortexbot

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/proxy"

	"github.com/boltdb/bolt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ilyaglow/go-cortex"
)

const (
	// defaultTLP is Green because indicators reaching telegram servers
	// TODO: think about making it configurable
	defaultTLP           = 1
	boltFileName         = "bolt.db"
	bucket               = "users"
	tgTokenEnvName       = "TGBOT_API_TOKEN"
	cortexURLEnvName     = "CORTEX_URL"
	cortexAPIKeyEnvName  = "CORTEX_API_KEY"
	socksURLEnvName      = "SOCKS5_URL"
	cortexBotPWEnvName   = "CORTEX_BOT_PASSWORD"
	cortexTimeoutEnvName = "CORTEX_TIMEOUT"
	debugEnvName         = "CORTEX_BOT_DEBUG"
)

var (
	pollTimeout   = 20 * time.Second
	cortexTimeout = 5 * time.Minute

	tgTokenEnvValue       = os.Getenv(tgTokenEnvName)
	cortexURLEnvValue     = os.Getenv(cortexURLEnvName)
	cortexAPIKeyEnvValue  = os.Getenv(cortexAPIKeyEnvName)
	cortexBotPWEnvValue   = os.Getenv(cortexBotPWEnvName)
	socksURLEnvValue      = os.Getenv(socksURLEnvName)
	cortexTimeoutEnvValue = os.Getenv(cortexTimeoutEnvName)
	debugEnvValue         = os.Getenv(debugEnvName)
)

// Client defines bot's abilities to interact with services
type Client struct {
	Bot         *tgbotapi.BotAPI
	Cortex      *cortex.Client
	Password    string
	DB          *bolt.DB
	UsersBucket string
	TLP         int
	Timeout     time.Duration
	Debug       bool
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
`, tgTokenEnvName, cortexURLEnvName, cortexAPIKeyEnvName, cortexBotPWEnvName, ex)
	os.Exit(1)
}

func (c *Client) log(v ...interface{}) {
	if c.Debug {
		log.Println(v...)
	}
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

	if tgTokenEnvValue == "" || cortexURLEnvValue == "" || cortexAPIKeyEnvValue == "" || cortexBotPWEnvValue == "" {
		usage()
	}

	if socksURLEnvValue != "" {
		surl, err := url.Parse(socksURLEnvValue)
		if err != nil {
			log.Fatal(err)
		}

		sc, err := socks5Client(surl)
		if err != nil {
			log.Fatal(err)
		}

		bot, err = tgbotapi.NewBotAPIWithClient(tgTokenEnvValue, sc)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		bot, err = tgbotapi.NewBotAPI(tgTokenEnvValue)
		if err != nil {
			log.Fatal(err)
		}
	}

	crtx, err := cortex.NewClient(cortexURLEnvValue, &cortex.ClientOpts{
		Auth: &cortex.APIAuth{
			APIKey: cortexAPIKeyEnvValue,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	db, err := bolt.Open(boltFileName, 0644, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a bucket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return errors.New("create users bucket failed")
		}
		return nil
	})

	timeout := cortexTimeout
	if cortexTimeoutEnvValue != "" {
		timeout, err = time.ParseDuration(cortexTimeoutEnvValue)
		if err != nil {
			log.Fatal(err)
		}
	}

	var debug bool
	debug, err = strconv.ParseBool(debugEnvValue)
	if err != nil {
		debug = false
	}

	return &Client{
		Bot:         bot,
		Cortex:      crtx,
		Password:    cortexBotPWEnvValue,
		DB:          db,
		UsersBucket: bucket,
		TLP:         defaultTLP,
		Timeout:     timeout,
		Debug:       debug,
	}
}
