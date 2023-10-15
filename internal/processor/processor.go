package processor

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type service struct {
	mp3DirPath string
	d          []Downloader
	u          []Uploader
}

func New(mp3Dir string, d []Downloader, u []Uploader) *service {
	return &service{
		mp3DirPath: mp3Dir,
		d:          d,
		u:          u,
	}
}

func (s *service) Process(ctx context.Context, videoURL string) (string, error) {
	title, err := getTitle(ctx, videoURL)
	if err != nil {
		return "", err
	}

	filename := title + ".mp3"

	outputPath := s.mp3DirPath + "/" + filename

	if err := s.download(ctx, videoURL, outputPath); err != nil {
		fmt.Println(err)
		return "", err
	}

	uploadRes := s.upload(ctx, videoURL, outputPath)
	fmt.Println(uploadRes)
	return uploadRes, nil
}

func (s *service) download(ctx context.Context, sourceURL, output string) error {
	for _, dd := range s.d {
		err := dd.Download(ctx, sourceURL, output)
		if err == nil {
			return nil
		}
		fmt.Printf("%s fail: %v\n", dd.Name(), err)
	}
	return errors.New("failed to download source")
}

func (s *service) upload(ctx context.Context, sourceURL, output string) string {
	var res []string
	for _, uu := range s.u {
		if err := uu.Upload(ctx, sourceURL, output); err != nil {
			res = append(res, fmt.Sprintf("%s upload fail: %v\n", uu.Name(), err))
			continue
		}
		res = append(res, fmt.Sprintf("%s upload success\n", uu.Name()))
	}

	return strings.Join(res, "\n")
}

func getTitle(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	node, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	if title, ok := traverse(node); ok {
		return title, nil
	}

	return "", errors.New("title not found")
}

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}
