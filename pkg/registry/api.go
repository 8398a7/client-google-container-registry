package registry

import (
	"context"
	"fmt"

	"golang.org/x/xerrors"
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
		return nil, xerrors.Errorf("failed to get registry token: %w", err)
	}

	path := "/tags/list"
	req, err := c.newRequest(ctx, "GET", path, nil, token.Token)
	if err != nil {
		return nil, xerrors.Errorf("failed to generate get images request(%s): %w", path, err)
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("failed to get images request(%s): %w", path, err)
	}

	var body ImageList
	if err := decodeBody(res, &body); err != nil {
		return nil, xerrors.Errorf("failed to decode body: %w", err)
	}
	return &body, nil
}

func (c *Client) GetTags(ctx context.Context, image string) (*ImageList, error) {
	token, err := c.getRegistryToken(ctx, "pull", image)
	if err != nil {
		return nil, xerrors.Errorf("failed to get registry token(%s): %w", image, err)
	}

	path := fmt.Sprintf("/%s/tags/list", image)
	req, err := c.newRequest(ctx, "GET", path, nil, token.Token)
	if err != nil {
		return nil, xerrors.Errorf("failed to generate get tags request(%s): %w", path, err)
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("failed to get tags request(%s): %w", path, err)
	}

	var body ImageList
	if err := decodeBody(res, &body); err != nil {
		return nil, xerrors.Errorf("failed to decode body: %w", err)
	}
	return &body, nil
}

func (c *Client) DeleteImage(ctx context.Context, image, tag string) (*RegistryError, error) {
	tags, err := c.GetTags(ctx, image)
	if err != nil {
		return nil, xerrors.Errorf("failed to get tags(%s): %w", image, err)
	}
	var deleteHashes []string
	for hash, manifest := range tags.Manifest {
		for _, t := range manifest.Tag {
			if tag == t {
				deleteHashes = append(deleteHashes, hash)
			}
		}
	}

	token, err := c.getRegistryToken(ctx, "push,pull", image)
	if err != nil {
		return nil, xerrors.Errorf("failed to get registry token(%s): %w", image, err)
	}

	path := fmt.Sprintf("/%s/manifests/%s", image, tag)
	req, err := c.newRequest(ctx, "DELETE", path, nil, token.Token)
	if err != nil {
		return nil, xerrors.Errorf("failed to generate delete request(%s): %w", path, err)
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("failed to delete request(%s): %w", path, err)
	}

	var body RegistryError
	if err := decodeBody(res, &body); err != nil {
		return nil, err
	}

	for _, hash := range deleteHashes {
		res, err := c.DeleteImageWithSha256(ctx, image, hash)
		if err != nil {
			return nil, xerrors.Errorf("failed to delete image with sha256: %w", err)
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
		return nil, xerrors.Errorf("failed to get registry token(%s): %w", image, err)
	}

	path := fmt.Sprintf("/%s/manifests/%s", image, hash)
	req, err := c.newRequest(ctx, "DELETE", path, nil, token.Token)
	if err != nil {
		return nil, xerrors.Errorf("failed to generate delete request(%s): %w", path, err)
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("failed to delete request(%s): %w", path, err)
	}

	var body RegistryError
	if err := decodeBody(res, &body); err != nil {
		return nil, xerrors.Errorf("failed to decode body: %w", err)
	}
	return &body, nil
}
