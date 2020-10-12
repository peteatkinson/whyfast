package utils

import (
	"encoding/binary"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"os"
	"time"
)

var buffer = make([][]byte, 0)

func PlaySound(s *discordgo.Session, guildId, channelId string) (err error) {
	vc, err := s.ChannelVoiceJoin(guildId, channelId, false, false)
	if err != nil {
		return err
	}

	time.Sleep(250 * time.Millisecond)

	fmt.Println("Bot is now speaking in channel...")
	vc.Speaking(true)

	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	fmt.Println("Bot no longer speaking in channel...")
	vc.Speaking(false)

	time.Sleep(250 * time.Millisecond)

	vc.Disconnect()
	return nil
}

func LoadSound() error {
	file, err := os.Open("./sounds/why_you_coming_fast.dca")

	if err != nil {
		fmt.Println("Error opening dca file:", err)
		return err
	}

	var opuslen int16

	for {
		// read opus frame length
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// if EOF, return
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file:", err)
			return err
		}

		if opuslen > 0 {
			// Read encoded pcm from dca
			InBuf := make([]byte, opuslen)
			err = binary.Read(file, binary.LittleEndian, &InBuf)
			// Not EOF
			if err != nil {
				fmt.Println("Error reading from dca file:", err)
				return err
			}

			buffer = append(buffer, InBuf)
		}
	}
}
