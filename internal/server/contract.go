package server

import "context"

type commandsProvider interface {
	NextCmd() (string, string, error)
	Close()
}

type processor interface {
	Process(ctx context.Context, videoURL string) (string, error)
}

type msgSender interface {
	SendMsg(text string) error
}
