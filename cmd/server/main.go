package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"getMP3/internal/downloader/youtubeDL"
	"getMP3/internal/downloader/ytDLP"
	"getMP3/internal/processor"
	"getMP3/internal/server"
	"getMP3/internal/shutdown"
	"getMP3/internal/uploader/tgBot"
)

func main() {
	var (
		ctx = context.Background()
		err error
	)
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

	closer := shutdown.New(time.Second)

	mp3Dir, found := os.LookupEnv("MP3_DIR")
	if !found {
		fmt.Println("MP3_DIR not found")
		return
	}

	bot, err := tgBot.New(botToken, int64(chatID))
	if err != nil {
		return
	}

	p := processor.New(mp3Dir,
		[]processor.Downloader{
			ytDLP.New(),
			youtubeDL.New(),
		},
		[]processor.Uploader{
			bot,
		},
	)

	s := server.New(bot, p, bot)
	closer.Add(s.Stop)

	go func(ctx context.Context) {
		closer.Wait(ctx)
		closer.CloseAll()
	}(ctx)

	s.Run(ctx)
}
