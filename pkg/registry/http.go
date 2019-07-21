package registry

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"path"

	"golang.org/x/xerrors"
)

func (c *Client) newRequest(ctx context.Context, method, spath string, body io.Reader, token string) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, xerrors.Errorf("failed to new request(method: %s, url: %s): %w", method, u.String(), err)
	}

	req = req.WithContext(ctx)

	req.Header.Add("Authorization", "Bearer "+token)

	return req, nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
