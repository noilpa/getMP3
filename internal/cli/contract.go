package cli

import "context"

type Processor interface {
	Process(ctx context.Context, videoURL string) (string, error)
}
