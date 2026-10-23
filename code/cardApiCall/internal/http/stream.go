package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/samber/do"
)

// StreamClient は HTTP ストリーミング取得のインターフェースです。
type StreamClient interface {
	GetStream(ctx context.Context, url string) (io.ReadCloser, error)
}

type defaultStreamClient struct {
	client *http.Client
}

// NewDefaultStreamClient はデフォルトの StreamClient を作成します。
func NewDefaultStreamClient() StreamClient {
	return &defaultStreamClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewStreamClientForDI は samber/do 用のプロバイダ関数です。
func NewStreamClientForDI(i *do.Injector) (StreamClient, error) {
	return NewDefaultStreamClient(), nil
}

// GetStream は指定 URL から画像ストリーム (io.ReadCloser) を取得します。
func (c *defaultStreamClient) GetStream(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected http status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
