package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyaglow/cortex-tgbot"
)

func main() {
	c, err := cortexbot.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	c.Bot.Debug = true

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
