package main

import (
	"flag"
	"fmt"
	bot "gptscratchadmin/internal/bot"
)

var OpenaiKey, DiscordKey string

func init() {
    flag.StringVar(&DiscordKey, "discordkey", "", "discord application token")
    flag.StringVar(&OpenaiKey, "openaikey", "", "openai api key")

    flag.Parse()
}

func main() {
    fmt.Println(OpenaiKey + " " + DiscordKey)
    bot.Run(DiscordKey, OpenaiKey)
}
