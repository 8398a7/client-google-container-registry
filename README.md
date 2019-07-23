# client-google-container-registry

## Usage

```bash
# require `roles/storage.admin` role
export GOOGLE_APPLICATION_CREDENTIALS=service-account.json
```

```go
func main() {
	ctx := context.Background()
	client, err := registry.NewClient("gcr.io", "your-gcp-project-id")
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
	res, err := client.DeleteImage(ctx, "image1", "tag"))
	if err != nil {
		// error handling
	}
	res, err := client.DeleteImageWithHash(ctx, "image1", "sha256:hash"))
	if err != nil {
		// error handling
	}
}
```
