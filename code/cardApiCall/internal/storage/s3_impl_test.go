package storage

import (
	"os"
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestS3StorageClient_GetDestinationURL(t *testing.T) {
	t.Run("uses cdnBaseURL when provided", func(t *testing.T) {
		client := NewS3StorageClient(nil, "my-bucket", "ap-northeast-1", "https://cdn.example.com")
		url := client.GetDestinationURL("cards/large/12345.jpg")
		assert.Equal(t, "https://cdn.example.com/cards/large/12345.jpg", url)

		// trailing slash in cdnBaseURL is handled
		clientTrailing := NewS3StorageClient(nil, "my-bucket", "ap-northeast-1", "https://cdn.example.com/")
		urlTrailing := clientTrailing.GetDestinationURL("/cards/large/12345.jpg")
		assert.Equal(t, "https://cdn.example.com/cards/large/12345.jpg", urlTrailing)
	})

	t.Run("uses standard s3 url when cdnBaseURL is empty", func(t *testing.T) {
		client := NewS3StorageClient(nil, "my-bucket", "ap-northeast-1", "")
		url := client.GetDestinationURL("cards/small/12345.jpg")
		assert.Equal(t, "https://my-bucket.s3.ap-northeast-1.amazonaws.com/cards/small/12345.jpg", url)
	})
}

func TestNewS3StorageClientForDI(t *testing.T) {
	// Clean env after test
	t.Cleanup(func() {
		os.Unsetenv("AWS_REGION")
		os.Unsetenv("CARD_IMAGE_BUCKET_NAME")
		os.Unsetenv("S3_BUCKET_NAME")
		os.Unsetenv("CARD_IMAGE_CDN_BASE_URL")
		os.Unsetenv("CDN_BASE_URL")
		os.Unsetenv("AWS_ENDPOINT_URL")
		os.Unsetenv("S3_ENDPOINT_URL")
		os.Unsetenv("AWS_S3_USE_PATH_STYLE")
	})

	t.Run("initializes storage client with environment variables", func(t *testing.T) {
		os.Setenv("AWS_REGION", "ap-northeast-1")
		os.Setenv("S3_BUCKET_NAME", "test-bucket")
		os.Setenv("CDN_BASE_URL", "http://localhost:9000/test-bucket")
		os.Setenv("AWS_ENDPOINT_URL", "http://minio:9000")
		os.Setenv("AWS_S3_USE_PATH_STYLE", "true")

		injector := do.New()
		do.Provide(injector, NewS3StorageClientForDI)

		client, err := do.Invoke[StorageClient](injector)
		require.NoError(t, err)
		require.NotNil(t, client)

		impl, ok := client.(*s3StorageClientImpl)
		require.True(t, ok)
		assert.Equal(t, "test-bucket", impl.bucketName)
		assert.Equal(t, "ap-northeast-1", impl.region)
		assert.Equal(t, "http://localhost:9000/test-bucket", impl.cdnBaseURL)
		assert.NotNil(t, impl.s3Client)
	})
}
