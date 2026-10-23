package storage

import (
	"context"
	"io"
)

// StorageClient は画像ストレージ操作の抽象インターフェースです。
type StorageClient interface {
	// ObjectExists は指定キーが S3 に存在するか確認する (HeadObject)
	ObjectExists(ctx context.Context, key string) (bool, error)

	// UploadStream は Reader から読み込んだバイナリを S3 にストリーミング転送する
	UploadStream(ctx context.Context, key string, body io.Reader, contentType string) (string, error)

	// GetDestinationURL はS3オブジェクトキーに対応する配信URL（CloudFrontまたはS3配信URL）を返す
	GetDestinationURL(key string) string
}
