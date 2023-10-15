package ytDLP

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type provider struct {
	binPath string
	binOpts string
}

func New() *provider {
	return &provider{
		binPath: "yt-dlp",
		binOpts: "--extract-audio --audio-format mp3  --output",
	}
}

func (p *provider) Download(ctx context.Context, sourceURL, output string) error {
	cmd := fmt.Sprintf("%s %s \"%s\" %s ", p.binPath, p.binOpts, output, sourceURL)
	fmt.Println(cmd)
	download := exec.CommandContext(ctx, "zsh", "-c", cmd)
	download.Stdout = os.Stdout
	download.Stderr = os.Stdin
	return download.Run()
}

func (p *provider) Name() string {
	return "yt-dlp"
}
