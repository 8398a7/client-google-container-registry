package registry

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

type registryToken struct {
	ExpiresIn float32   `json:"expires_in"`
	IssuedAt  time.Time `json:"issued_at"`
	Token     string    `json:"token"`
}

func (c *Client) getRegistryToken(ctx context.Context, image string) (*registryToken, error) {
	u, err := url.ParseRequestURI("https://gcr.io/v2/token")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	if image == "" {
		q.Add("scope", fmt.Sprintf("repository:%s:pull", c.repo))
	} else {
		q.Add("scope", fmt.Sprintf("repository:%s/%s:pull", c.repo, image))
	}
	q.Add("service", c.host)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	uEnc := base64.URLEncoding.EncodeToString([]byte("_token:" + c.token))
	req.Header.Add("Host", c.host)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+uEnc)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	var token *registryToken
	if err := decodeBody(res, &token); err != nil {
		return nil, err
	}

	return token, nil
}

func getAccessToken(ctx context.Context) (string, error) {
	out, err := exec.Command("gcloud", "auth", "print-access-token").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
