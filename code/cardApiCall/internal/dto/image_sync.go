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
	NeuronID       int64
	OcgApiImageID  int64
	Variant        ImageVariant // Small, Large, Cropped
	SourceURL      string
	S3Key          string
	DestinationURL string
}

// SyncImagesInput はLambdaの入力イベントDTOです。
type SyncImagesInput struct {
	NeuronCardIds []int64 `json:"neuronCardIds"`
}

// SyncError は同期処理中のエラー詳細DTOです。
type SyncError struct {
	NeuronCardID  int64  `json:"neuronCardId"`
	OcgApiImageID int64  `json:"ocgApiImageId"`
	Message       string `json:"message"`
}

// SyncImagesOutput はLambdaの出力イベントDTOです。
type SyncImagesOutput struct {
	TotalProcessedCards     int         `json:"totalProcessedCards"`
	SuccessfulNeuronCardIds []int64     `json:"successfulNeuronCardIds"`
	FailedNeuronCardIds     []int64     `json:"failedNeuronCardIds"`
	SyncedImageCount        int         `json:"syncedImageCount"`
	SkippedImageCount       int         `json:"skippedImageCount"`
	Errors                  []SyncError `json:"errors"`
}

