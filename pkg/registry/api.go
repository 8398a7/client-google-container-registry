package registry

import (
	"context"
	"fmt"
	"strings"
)

type Manifest struct {
	ImageSizeBytes string   `json:"imageSizeBytes"`
	LayerID        string   `json:"layerId"`
	MediaType      string   `json:"mediaType"`
	Tag            []string `json:"tag"`
	TimeCreatedMs  string   `json:"timeCreatedMs"`
	TimeUploadedMs string   `json:"timeUploadedMs"`
}

type ImageList struct {
	Child    []string            `json:"child"`
	Manifest map[string]Manifest `json:"manifest"`
	Name     string              `json:"name"`
	Tags     []string            `json:"tags"`
}

type RegistryError struct {
	Errors []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

func (c *Client) GetImages(ctx context.Context) (*ImageList, error) {
	token, err := c.getRegistryToken(ctx, "pull", "")
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, "GET", "/tags/list", nil, token.Token)
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	var body ImageList
	if err := decodeBody(res, &body); err != nil {
		return nil, err
	}
	return &body, nil
}

func (c *Client) GetTags(ctx context.Context, image string) (*ImageList, error) {
	token, err := c.getRegistryToken(ctx, "pull", image)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/%s/tags/list", image), nil, token.Token)
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	var body ImageList
	if err := decodeBody(res, &body); err != nil {
		return nil, err
	}
	return &body, nil
}

func (c *Client) DeleteImage(ctx context.Context, image, tag string) (*RegistryError, error) {
	tags, err := c.GetTags(ctx, image)
	if err != nil {
		return nil, err
	}
	var deleteHashes []string
	for sha256, manifest := range tags.Manifest {
		for _, t := range manifest.Tag {
			if tag == t {
				hash := strings.Split(sha256, ":")[1]
				deleteHashes = append(deleteHashes, hash)
			}
		}
	}

	token, err := c.getRegistryToken(ctx, "push,pull", image)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, "DELETE", fmt.Sprintf("/%s/manifests/%s", image, tag), nil, token.Token)
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	var body RegistryError
	if err := decodeBody(res, &body); err != nil {
		return nil, err
	}

	for _, hash := range deleteHashes {
		res, err := c.DeleteImageWithSha256(ctx, image, hash)
		if err != nil {
			return nil, err
		}
		if len(res.Errors) > 0 {
			return res, nil
		}
	}

	return &body, nil
}

func (c *Client) DeleteImageWithSha256(ctx context.Context, image, hash string) (*RegistryError, error) {
	token, err := c.getRegistryToken(ctx, "push,pull", image)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, "DELETE", fmt.Sprintf("/%s/manifests/sha256:%s", image, hash), nil, token.Token)
	if err != nil {
		return nil, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	var body RegistryError
	if err := decodeBody(res, &body); err != nil {
		return nil, err
	}
	return &body, nil
}
