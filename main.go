package main

import (
	"fmt"
	"flag"
	"strings"
	"syscall"
	"os/signal"
	"os"
	"log"
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

	bot.AddHandler(messageCreate)

	if err := bot.Open(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bot is now running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill) 
	<-sc

	// close down the bot session
	bot.Close()

}

func voiceStateUpdated(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	if v.UserID == s.State.User.ID {
		return
	}

	var hasRole bool
	hasRole = checkRoleStatus(s, v.GuildID, v.UserID, "WhyFast")
	if !hasRole {
		fmt.Println(hasRole)
		return
	}
	return
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
		fmt.Println(s)
		fmt.Println(m)


		c, err := s.State.Channel(m.ChannelID)

		if err != nil {
			return
		}

		guild, err := s.State.Guild(c.GuildID)

		if err != nil {
			return
		}


		if len(guild.VoiceStates) == 0 {
			fmt.Println("No voice channels to connect to")
		}

		for _, vs := range guild.VoiceStates {
			if vs.UserID == m.Author.ID {
				fmt.Println("Playing sound")
			}
			return
		}
	}
}

func checkRoleStatus(s * discordgo.Session, guildId string, userId string, role string) bool {
	member, err := s.GuildMember(guildId, userId)

	if err != nil {
		return false
	}

	for r := range member.Roles {
		if member.Roles[r] == role {
			return true
		} 
	}

	return false
}

func guildLeave() {

}


