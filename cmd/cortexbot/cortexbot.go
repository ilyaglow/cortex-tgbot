package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ilyaglow/cortex-tgbot"
	cortex "github.com/ilyaglow/go-cortex"
)

func main() {
	c, err := cortexbot.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	c.Bot.Debug = true

	crtx, err := cortex.NewClient(os.Getenv("CORTEX_URL"), &cortex.ClientOpts{
		Auth: &cortex.APIAuth{
			APIKey: os.Getenv("CORTEX_API_KEY"),
		},
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	c.Cortex = crtx

	ctx, cancel := context.WithCancel(context.Background())
	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-ctrlc
		log.Println("main: got sigint")
		cancel()
		return
	}()

	log.Fatal(c.Run(ctx))
}
