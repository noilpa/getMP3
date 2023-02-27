package main

import (
	"fmt"
	"os"
	"strconv"

	"getMP3/internal/cli"
	"getMP3/internal/downloader/youtubeDL"
	"getMP3/internal/downloader/ytDLP"
	"getMP3/internal/uploader/tgBot"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
	}()

	botToken, found := os.LookupEnv("BOT_TOKEN")
	if !found {
		fmt.Println("BOT_TOKEN not found")
		return
	}

	chatIdStr, found := os.LookupEnv("CHAT_ID")
	if !found {
		fmt.Println("CHAT_ID not found")
		return
	}
	chatID, err := strconv.Atoi(chatIdStr)
	if err != nil {
		return
	}

	mp3Dir, found := os.LookupEnv("MP3_DIR")
	if !found {
		fmt.Println("MP3_DIR not found")
		return
	}

	bot, err := tgBot.New(botToken, int64(chatID))
	if err != nil {
		return
	}

	cli.New(
		mp3Dir,
		[]cli.Downloader{
			ytDLP.New(),
			youtubeDL.New(),
		},
		[]cli.Uploader{
			bot,
		},
	).Run()
}
