package server

import (
	"context"
	"fmt"
	"strings"
)

type server struct {
	cp commandsProvider
	p  processor
	ms msgSender
}

func New(cp commandsProvider, p processor, ms msgSender) *server {
	return &server{
		cp: cp,
		p:  p,
		ms: ms,
	}
}

func (s *server) Run(ctx context.Context) {
	defer fmt.Println("server stopped")
	fmt.Println("server started")
	for {
		cmd, args, err := s.cp.NextCmd()
		if err != nil {
			// log err
			return
		}

		fmt.Printf("command received: %s %s\n", cmd, args)

		switch cmd {
		case "mp3":
			s.processMP3(ctx, args)
		default:
			s.processDefault(ctx, cmd, args)
		}
	}
}

func (s *server) Stop() error {
	s.cp.Close()
	return nil
}

func (s *server) processMP3(ctx context.Context, args string) {
	videoURL := strings.Split(args, " ")[0]
	_, err := s.p.Process(ctx, videoURL)
	if err != nil {
		fmt.Printf("processMP3 failed to process: %v", err)
		if err2 := s.ms.SendMsg(fmt.Sprintf("failed to process: %v", err)); err2 != nil {
			fmt.Printf("processMP3 send message error: %v\n", err2)
		}
		return
	}

	return
}

func (s *server) processDefault(ctx context.Context, cmd, args string) {
	err := s.ms.SendMsg(fmt.Sprintf("unknown command: %s %s", cmd, args))
	if err != nil {
		fmt.Printf("processDefault send message error: %v\n", err)
	}
}
