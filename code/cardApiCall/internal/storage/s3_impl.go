package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"atomisu.com/ocg-statics/cardApiCall/internal/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithy "github.com/aws/smithy-go"
	"github.com/samber/do"
)

type s3StorageClientImpl struct {
	s3Client        *s3.Client
	transferManager *transfermanager.Client
	bucketName      string
	region          string
	cdnBaseURL      string
}

// NewS3StorageClient は S3 クライアントのインスタンスを生成します。
func NewS3StorageClient(s3Client *s3.Client, bucketName string, region string, cdnBaseURL string) StorageClient {
	return &s3StorageClientImpl{
		s3Client:        s3Client,
		transferManager: transfermanager.New(s3Client),
		bucketName:      bucketName,
		region:          region,
		cdnBaseURL:      strings.TrimSuffix(cdnBaseURL, "/"),
	}
}

// NewS3StorageClientForDI は samber/do コンテナ用のプロバイダ関数です。
func NewS3StorageClientForDI(i *do.Injector) (StorageClient, error) {
	ctx := context.Background()
	region := env.GetAsString("AWS_REGION", "ap-northeast-1")
	bucketName := env.GetAsString("CARD_IMAGE_BUCKET_NAME", env.GetAsString("S3_BUCKET_NAME", ""))
	cdnBaseURL := env.GetAsString("CARD_IMAGE_CDN_BASE_URL", env.GetAsString("CDN_BASE_URL", ""))
	endpointURL := env.GetAsString("AWS_ENDPOINT_URL", env.GetAsString("S3_ENDPOINT_URL", ""))
	usePathStyle := env.GetAsString("AWS_S3_USE_PATH_STYLE", env.GetAsString("S3_USE_PATH_STYLE", "false")) == "true"

	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load aws sdk config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpointURL != "" {
			o.BaseEndpoint = aws.String(endpointURL)
			o.UsePathStyle = true
		}
		if usePathStyle {
			o.UsePathStyle = true
		}
	})
	return NewS3StorageClient(s3Client, bucketName, region, cdnBaseURL), nil
}

// ObjectExists は指定されたS3キーが存在するか確認します。
func (s *s3StorageClientImpl) ObjectExists(ctx context.Context, key string) (bool, error) {
	_, err := s.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	})
	if err == nil {
		return true, nil
	}

	var notFound *types.NotFound
	if errors.As(err, &notFound) {
		return false, nil
	}

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		if apiErr.ErrorCode() == "NotFound" || apiErr.ErrorCode() == "NoSuchKey" {
			return false, nil
		}
	}

	if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "NotFound") {
		return false, nil
	}

	return false, fmt.Errorf("failed to head object for key %s: %w", key, err)
}

// UploadStream は Reader から読み込んだストリームを S3 にアップロードします。
func (s *s3StorageClientImpl) UploadStream(ctx context.Context, key string, body io.Reader, contentType string) (string, error) {
	input := &transfermanager.UploadObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
		Body:   body,
	}
	if contentType != "" {
		input.ContentType = &contentType
	}

	result, err := s.transferManager.UploadObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload stream to s3 key %s: %w", key, err)
	}

	if result.Location != nil && *result.Location != "" {
		return *result.Location, nil
	}
	return s.GetDestinationURL(key), nil
}

// GetDestinationURL は S3 オブジェクトキーに対するアクセス URL を返します。
func (s *s3StorageClientImpl) GetDestinationURL(key string) string {
	cleanKey := strings.TrimPrefix(key, "/")
	if s.cdnBaseURL != "" {
		return fmt.Sprintf("%s/%s", s.cdnBaseURL, cleanKey)
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, cleanKey)
}
