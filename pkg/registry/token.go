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

	"golang.org/x/xerrors"
)

type registryToken struct {
	ExpiresIn float32   `json:"expires_in"`
	IssuedAt  time.Time `json:"issued_at"`
	Token     string    `json:"token"`
}

const tokenUrl = "https://gcr.io/v2/token"

func (c *Client) getRegistryToken(ctx context.Context, grant, image string) (*registryToken, error) {
	u, err := url.ParseRequestURI(tokenUrl)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse request uri(%s): %w", tokenUrl, err)
	}
	q := u.Query()
	if image == "" {
		q.Add("scope", fmt.Sprintf("repository:%s:%s", c.repo, grant))
	} else {
		q.Add("scope", fmt.Sprintf("repository:%s/%s:%s", c.repo, image, grant))
	}
	q.Add("service", c.host)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, xerrors.Errorf("failed to generate request for get registry token: %w", err)
	}

	req = req.WithContext(ctx)

	uEnc := base64.URLEncoding.EncodeToString([]byte("_token:" + c.token))
	req.Header.Add("Host", c.host)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+uEnc)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("failed to request for get registry token: %w", err)
	}
	var token *registryToken
	if err := decodeBody(res, &token); err != nil {
		return nil, xerrors.Errorf("failed to decode body: %w", err)
	}

	return token, nil
}

func getAccessToken(ctx context.Context) (string, error) {
	out, err := exec.Command("gcloud", "auth", "print-access-token").Output()
	if err != nil {
		return "", xerrors.Errorf("failed to execute `gcloud auth print-access-token`: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
