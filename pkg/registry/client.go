package registry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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
		return nil, err
	}
	reqUrl, err := url.ParseRequestURI(fmt.Sprintf("https://%s/v2/%s", host, repo))
	if err != nil {
		return nil, err
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
