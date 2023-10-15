package cli

import (
	"context"
	"fmt"
	"os"
)

type service struct {
	p Processor
}

func New(p Processor) *service {
	return &service{
		p: p,
	}
}

func (s *service) Run() {
	if len(os.Args) != 2 {
		fmt.Println("source url not found")
		return
	}

	ctx := context.Background()
	videoURL := os.Args[1]

	res, err := s.p.Process(ctx, videoURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}
