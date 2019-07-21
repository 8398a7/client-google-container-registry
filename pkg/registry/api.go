package registry

import (
	"context"
	"fmt"
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

func (c *Client) GetImages(ctx context.Context) (*ImageList, error) {
	req, err := c.newRequest(ctx, "GET", "/tags/list", nil, "")
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
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/%s/tags/list", image), nil, image)
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
