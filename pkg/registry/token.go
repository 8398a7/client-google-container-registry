package registry

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jws"
	"golang.org/x/xerrors"
)

const audience = "https://www.googleapis.com/oauth2/v4/token"
const scope = "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/cloud-platform"
const grantType = "urn:ietf:params:oauth:grant-type:jwt-bearer"

type registryToken struct {
	ExpiresIn float32   `json:"expires_in"`
	IssuedAt  time.Time `json:"issued_at"`
	Token     string    `json:"token"`
}

func (c *Client) getRegistryToken(ctx context.Context, grant, image string) (*registryToken, error) {
	u, err := url.ParseRequestURI(c.getTokenUrl())
	if err != nil {
		return nil, xerrors.Errorf("failed to parse request uri(%s): %w", c.getTokenUrl(), err)
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

func (c *Client) getTokenUrl() string {
	return fmt.Sprintf("https://%s/v2/token", c.host)
}

func generateJWT(saKeyfile string) (string, error) {
	sa, err := ioutil.ReadFile(saKeyfile)
	if err != nil {
		return "", xerrors.Errorf("Could not read service account file: %w", err)
	}
	conf, err := google.JWTConfigFromJSON(sa)
	if err != nil {
		return "", xerrors.Errorf("Could not parse service account JSON: %w", err)
	}
	block, _ := pem.Decode(conf.PrivateKey)
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", xerrors.Errorf("private key parse error: %w", err)
	}
	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return "", xerrors.New("private key failed rsa.PrivateKey type assertion")
	}
	now := time.Now().Unix()

	jwt := &jws.ClaimSet{
		Iat:   now,
		Exp:   now + 3600,
		Iss:   conf.Email,
		Aud:   audience,
		Scope: scope,
	}
	jwsHeader := &jws.Header{
		Algorithm: "RS256",
		Typ:       "JWT",
	}
	token, err := jws.Encode(jwsHeader, jwt, rsaKey)
	if err != nil {
		return "", xerrors.Errorf("failed to encode jws: %w", err)
	}

	return generateAccessToken(token)
}

func generateAccessToken(token string) (string, error) {
	client := new(http.Client)

	values := url.Values{}
	values.Add("grant_type", grantType)
	values.Add("assertion", token)

	req, err := http.NewRequest("POST", audience, strings.NewReader(values.Encode()))
	if err != nil {
		return "", xerrors.Errorf("failed to generate http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", xerrors.Errorf("failed to request token: %w", err)
	}

	var body struct {
		AccessToken string `json:"access_token"`
	}

	if err := decodeBody(resp, &body); err != nil {
		return "", xerrors.Errorf("failed to decode body: %w", err)
	}

	return body.AccessToken, nil
}
