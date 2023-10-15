package tgBot

import (
	"context"
	"errors"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	once *sync.Once

	ErrEventsIsClosed = errors.New("events channel is closed")
)

type provider struct {
	bot    *tgbotapi.BotAPI
	chatID int64
	events tgbotapi.UpdatesChannel
}

func init() {
	once = new(sync.Once)
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

func (p *provider) NextCmd() (string, string, error) {
	once.Do(func() {
		p.events = p.bot.GetUpdatesChan(tgbotapi.UpdateConfig{})
	})

	for e := range p.events {
		if e.Message == nil {
			continue
		}

		if !e.Message.IsCommand() {
			continue
		}

		return e.Message.Command(), e.Message.CommandArguments(), nil
	}

	return "", "", ErrEventsIsClosed
}

func (p *provider) Close() {
	p.bot.StopReceivingUpdates()
}

func (p *provider) SendMsg(text string) error {
	msg := tgbotapi.NewMessage(p.chatID, text)
	_, err := p.bot.Send(msg)

	return err
}

func (p *provider) Name() string {
	return "telegram-bot"
}
