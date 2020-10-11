package whyfast

import (
	"fmt"
	"time"
	"os"
	"encoding/binary"
	"io"
	"github.com/bwmarrin/discordgo"
)

var buffer = make([][]byte, 0)

func playSound(s *discordgo.Session, guildId, channelId string) (err error) {
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

func loadSound(soundEffect *string) error {
	file, err := os.Open("")

	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16
	for {
		err := binary.Read(file, binary.LittleEndian, &opuslen)

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file : ", err)
		}

		inBuffer := make([]byte, opuslen)

		err = binary.Read(file, binary.LittleEndian, &inBuffer)

		if err != nil {
			fmt.Println("Error reading from dca file : ", err)
			return err
		}

		buffer = append(buffer, inBuffer)
	}

}