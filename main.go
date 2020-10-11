package main

import (
	"fmt"
	"flag"
	"github.com/bwmarrin/discordgo"
)

var token string

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}


func main() {
	if token == "" {
		fmt.Println("No token provided. Please run: whyfast -t <bot token>")
		return
	}

	bot, err := discordgo.New(fmt.Sprintf("Bot %s", token))

	if err != nil {
		fmt.Println("Error creating Bot session", err)
	}
}


