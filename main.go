package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/peteratkinson/whyfast/utils"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	bot.AddHandler(messageCreate)

	if err := bot.Open(); err != nil {
		log.Fatal(err)
	}

	err = utils.LoadSound()
	if err != nil {
		fmt.Println("Could not load .dca file")
	}

	fmt.Println("Bot is now running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// close down the bot session
	bot.Close()

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Fields(m.Content)

	if len(args) == 0 {
		return
	}

	if args[0] != "!whyfast" {
		return
	} else {
		c, err := s.State.Channel(m.ChannelID)

		if err != nil {
			return
		}

		guild, err := s.State.Guild(c.GuildID)

		if err != nil || guild == nil {
			return
		}

		// find guild the user is currently in
		g, err := guildFromMessage(m, s)

		if g != nil {
			// find the channel the user is currently in
			cId, err := findUserChannel(m, s)

			// if the user is in a voicer channel then play the sound in that channel
			if err == nil {
				err = utils.PlaySound(s, g.ID, cId)
			}
		}
	}
}

func findUserChannel(m *discordgo.MessageCreate, s *discordgo.Session) (string, error) {
	// find the voice channel the user is in
	g, err := guildFromMessage(m, s)
	if err != nil {
		return "", err
	}

	for _, vstate := range g.VoiceStates {
		if vstate.UserID == m.Author.ID {
			return vstate.ChannelID, nil
		}
	}
	return "", fmt.Errorf("User is not in a voice channel")
}

func guildFromMessage(m *discordgo.MessageCreate, s *discordgo.Session) (*discordgo.Guild, error) {
	c, err := s.State.Channel(m.Message.ChannelID)
	if err != nil {
		return nil, err
	}
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		return nil, err
	}
	return g, nil
}
