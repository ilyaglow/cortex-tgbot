package main

import (
	"github.com/ilyaglow/cortex-tgbot"
)

func main() {
	c := cortexbot.NewClient()
	c.Bot.Debug = true
	c.Run()
}
