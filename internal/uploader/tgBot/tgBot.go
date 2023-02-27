package tgBot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type provider struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func New(token string, chatID int64) (*provider, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &provider{
		bot:    bot,
		chatID: chatID,
	}, nil
}

func (p *provider) Upload(ctx context.Context, sourceURL, outputPath string) error {
	ac := tgbotapi.NewAudio(p.chatID, tgbotapi.FilePath(outputPath))
	ac.Caption = sourceURL

	out, err := p.bot.Send(ac)
	if err != nil {
		fmt.Println(out)
		return err
	}

	return nil
}

func (p *provider) Name() string {
	return "telegram-bot"
}
