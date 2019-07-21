# client-google-container-registry

## Usage

```go
func main() {
	ctx := context.Background()
	client, err := registry.NewClient(ctx, "gcr.io", "your-gcp-project-id")
	if err != nil {
		// error handling
	}
	images, err := client.GetImages(ctx)
	if err != nil {
		// error handling
	}
	tags, err := client.GetTags(ctx, "image1")
	if err != nil {
		// error handling
	}
	res, err := client.DeleteImage(ctx, "image", "tag"))
	if err != nil {
		// error handling
	}
	res, err := client.DeleteImageWithHash(ctx, "image", "hash"))
	if err != nil {
		// error handling
	}
}
```
