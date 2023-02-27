package cli

import "context"

type Downloader interface {
	Download(ctx context.Context, sourceURL, output string) error
	Name() string
}

type Uploader interface {
	Upload(ctx context.Context, sourceURL, outputPath string) error
	Name() string
}
