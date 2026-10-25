package main

import (
	"context"
	"encoding/json"
	"fmt"

	"atomisu.com/ocg-statics/cardApiCall/cmd/aws_lambda/syncStorageImageWithDB/app"
	"atomisu.com/ocg-statics/cardApiCall/internal/dto"
	"atomisu.com/ocg-statics/cardApiCall/internal/usecase/syncimage"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/samber/do"
)

var globalInjector = app.SetupDIContainer()

// HandleRequest は syncStorageImageWithDB Lambda 関数のリクエストハンドラーです。
func HandleRequest(ctx context.Context, event json.RawMessage) (*dto.SyncImagesOutput, error) {
	var input dto.SyncImagesInput
	if err := json.Unmarshal(event, &input); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request event: %w", err)
	}

	if len(input.NeuronCardIds) == 0 {
		return nil, fmt.Errorf("neuronCardIds must not be empty")
	}

	syncUseCase := do.MustInvoke[syncimage.SyncStorageImageUseCase](globalInjector)
	output, err := syncUseCase.SyncStorageImagesByNeuronIDs(ctx, input.NeuronCardIds)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func main() {
	defer globalInjector.Shutdown()
	lambda.Start(HandleRequest)
}
