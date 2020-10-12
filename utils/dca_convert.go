package utils

import (
	"fmt"
	"github.com/jonas747/dca"
	"io"
	"log"
	"os"
)

func ConvertMp3Audio(file string) {
	options := dca.StdEncodeOptions
	// default options (use 120 bit rate for any mp3 downloaded from Youtube)
	options.RawOutput = true
	options.Bitrate = 120

	s, err := dca.EncodeFile(fmt.Sprintf("./sounds/%s.mp3", file), options)

	defer s.Cleanup()

	if err != nil {
		log.Fatal("Failed creating an encoding session: ", err)
	}

	output, err := os.Create("./sounds/test.dca")
	if err != nil {
		fmt.Println("Error creating .dca file to disk")
	}

	io.Copy(output, s)
}
