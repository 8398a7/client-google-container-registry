package registry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/xerrors"
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
	token      string
	host       string
	repo       string
}

func NewClient(ctx context.Context, host, repo string) (*Client, error) {
	token, err := getAccessToken(ctx)
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
