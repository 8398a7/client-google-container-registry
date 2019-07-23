package registry

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/xerrors"
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
	token      string
	host       string
	repo       string
}

func NewClient(host, repo string) (*Client, error) {
	keyFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if keyFile == "" {
		return nil, xerrors.New("GOOGLE_APPLICATION_CREDENTIALS is not specified.")
	}

	token, err := generateJWT(keyFile)
	if err != nil {
		return nil, xerrors.Errorf("failed to get access token: %w", err)
	}
	uri := fmt.Sprintf("https://%s/v2/%s", host, repo)
	reqUrl, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse request uri(%s): %w", uri, err)
	}

	client := new(http.Client)

	return &Client{
		URL:        reqUrl,
		HTTPClient: client,
		token:      token,
		host:       host,
		repo:       repo,
	}, nil
}
