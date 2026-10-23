package dto

// ImageVariant はカード画像のバリエーションを表す型です。
type ImageVariant string

const (
	ImageVariantSmall   ImageVariant = "small"
	ImageVariantLarge   ImageVariant = "large"
	ImageVariantCropped ImageVariant = "cropped"
)

// ImageSyncTaskDTO は画像同期タスクの事前計算用DTOです。
type ImageSyncTaskDTO struct {
	CardImageID    int64
	OcgApiCardID   int64
	OcgApiImageID  int64
	Variant        ImageVariant // Small, Large, Cropped
	SourceURL      string
	S3Key          string
	DestinationURL string
}

// SyncImagesInput はLambdaの入力イベントDTOです。
type SyncImagesInput struct {
	OcgCardIds []int64 `json:"ocgCardIds"`
}

// SyncError は同期処理中のエラー詳細DTOです。
type SyncError struct {
	OcgCardId     int64  `json:"ocgCardId"`
	OcgApiImageId int64  `json:"ocgApiImageId"`
	Message       string `json:"message"`
}

// SyncImagesOutput はLambdaの出力イベントDTOです。
type SyncImagesOutput struct {
	TotalProcessedCards  int         `json:"totalProcessedCards"`
	SuccessfulOcgCardIds []int64     `json:"successfulOcgCardIds"`
	FailedOcgCardIds     []int64     `json:"failedOcgCardIds"`
	SyncedImageCount     int         `json:"syncedImageCount"`
	SkippedImageCount    int         `json:"skippedImageCount"`
	Errors               []SyncError `json:"errors"`
}
