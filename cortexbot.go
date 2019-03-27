package cortexbot

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/proxy"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	cortex "github.com/ilyaglow/go-cortex/v3"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

const (
	sqliteFileName       = "cortexbot.db"
	tgTokenEnvName       = "TGBOT_API_TOKEN"
	cortexURLEnvName     = "CORTEX_URL"
	cortexAPIKeyEnvName  = "CORTEX_API_KEY"
	socksURLEnvName      = "SOCKS5_URL"
	cortexBotPWEnvName   = "CORTEX_BOT_PASSWORD"
	cortexTimeoutEnvName = "CORTEX_TIMEOUT"
	debugEnvName         = "CORTEX_BOT_DEBUG"
)

var (
	// defaultTLP is Green because indicators reach telegram servers.
	// Same for defaultPAP.
	// TODO: think about making it configurable.
	defaultTLP = cortex.TLPGreen
	defaultPAP = cortex.PAPGreen

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

// Client defines bot's abilities to interact with services.
// Deprecated: use Cortexbot instead.
type Client Cortexbot

// Cortexbot defines bot's abilities to interact with services.
type Cortexbot struct {
	Bot      *tgbotapi.BotAPI
	Cortex   *cortex.Client
	Password string
	DB       *sql.DB
	dbpath   string
	TLP      cortex.TLP
	PAP      cortex.PAP
	Timeout  time.Duration
	Debug    bool
}

func (c *Cortexbot) log(v ...interface{}) {
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

// NewClient bootstraps the Client struct from env variables.
// Deprecated: Use NewFromEnv, or New method instead.
func NewClient() *Client {
	client, err := NewFromEnv()
	if err != nil {
		panic(err)
	}
	return ((*Client)(client))
}

// NewFromEnv bootstraps Cortexbot from environment variables.
func NewFromEnv() (*Cortexbot, error) {
	var (
		bot *tgbotapi.BotAPI
		err error
	)

	if tgTokenEnvValue == "" || cortexURLEnvValue == "" || cortexAPIKeyEnvValue == "" || cortexBotPWEnvValue == "" {
		return nil, fmt.Errorf(
			"not enough parameters: check that environment variables %s, %s, %s and %s are set",
			tgTokenEnvName,
			cortexURLEnvName,
			cortexAPIKeyEnvName,
			cortexBotPWEnvName,
		)
	}

	if socksURLEnvValue != "" {
		surl, err := url.Parse(socksURLEnvValue)
		if err != nil {
			return nil, err
		}

		sc, err := socks5Client(surl)
		if err != nil {
			return nil, err
		}

		bot, err = tgbotapi.NewBotAPIWithClient(tgTokenEnvValue, sc)
		if err != nil {
			return nil, err
		}
	} else {
		bot, err = tgbotapi.NewBotAPI(tgTokenEnvValue)
		if err != nil {
			return nil, err
		}
	}

	crtx, err := cortex.NewClient(cortexURLEnvValue, &cortex.ClientOpts{
		Auth: &cortex.APIAuth{
			APIKey: cortexAPIKeyEnvValue,
		},
	})
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", sqliteFileName)
	if err != nil {
		return nil, err
	}

	timeout := cortexTimeout
	if cortexTimeoutEnvValue != "" {
		timeout, err = time.ParseDuration(cortexTimeoutEnvValue)
		if err != nil {
			return nil, err
		}
	}

	var debug bool
	debug, err = strconv.ParseBool(debugEnvValue)
	if err != nil {
		debug = false
	}

	c := &Cortexbot{
		Bot:      bot,
		Cortex:   crtx,
		Password: cortexBotPWEnvValue,
		DB:       db,
		TLP:      defaultTLP,
		PAP:      defaultPAP,
		Timeout:  timeout,
		Debug:    debug,
	}
	err = c.createUsersTbl()
	if err != nil {
		return nil, err
	}

	return c, nil
}

// SetupChatbot constructs a client to the messenger.
func SetupChatbot(token string, client *http.Client) func(*Cortexbot) error {
	return func(c *Cortexbot) error {
		if client == nil {
			bot, err := tgbotapi.NewBotAPI(token)
			if err != nil {
				return err
			}
			c.Bot = bot
			return nil
		}

		bot, err := tgbotapi.NewBotAPIWithClient(token, client)
		if err != nil {
			return err
		}
		c.Bot = bot
		return nil
	}
}

// SetCortex sets a cortex.Client.
func SetCortex(client *cortex.Client) func(*Cortexbot) error {
	return func(c *Cortexbot) error {
		c.Cortex = client
		return nil
	}
}

// SetDBPath sets a sqlitedb path.
func SetDBPath(n string) func(*Cortexbot) error {
	return func(c *Cortexbot) error {
		c.dbpath = n
		return nil
	}
}

// SetTLP sets TLP as an option.
func SetTLP(tlp cortex.TLP) func(*Cortexbot) error {
	return func(c *Cortexbot) error {
		c.TLP = tlp
		return nil
	}
}

// SetPAP sets PAP as an option.
func SetPAP(pap cortex.PAP) func(*Cortexbot) error {
	return func(c *Cortexbot) error {
		c.PAP = pap
		return nil
	}
}

// SetCortexTimeout will set the timeout that the client will wait for a
// response from Cortex at most.
func SetCortexTimeout(t time.Duration) func(*Cortexbot) error {
	return func(c *Cortexbot) error {
		c.Timeout = t
		return nil
	}
}

// Debug sets debug mode for the Cortexbot.
func Debug() func(*Cortexbot) error {
	return func(c *Cortexbot) error {
		c.Debug = true
		return nil
	}
}

// New bootstraps cortexbot configuration.
func New(opts ...func(*Cortexbot) error) (*Cortexbot, error) {
	cortexbot := &Cortexbot{
		dbpath:  sqliteFileName,
		TLP:     defaultTLP,
		PAP:     defaultPAP,
		Timeout: cortexTimeout,
		Debug:   false,
	}

	for _, option := range opts {
		err := option(cortexbot)
		if err != nil {
			return nil, err
		}
	}

	return cortexbot, nil
}
