package main

import (
	"context"
	"encoding/json"
	"fmt"

	"atomisu.com/ocg-statics/cardApiCall/cmd/aws_lambda/syncStorageImageWithDB/app"
	"atomisu.com/ocg-statics/cardApiCall/internal/usecase/syncimage"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/samber/do"
)

type Input struct {
	StartId int64 `json:"startId"`
	Delta   int64 `json:"delta"`
}

type Output struct {
	FailedIds []int64 `json:"failedIds"`
}

var globalInjector = app.SetupDIContainer()

func HandleRequest(ctx context.Context, event json.RawMessage) (Output, error) {
	var input Input
	if err := json.Unmarshal(event, &input); err != nil {
		return Output{}, err
	}

	if input.Delta <= 0 {
		return Output{}, fmt.Errorf("delta must be > 0")
	}

	syncUseCase := do.MustInvoke[syncimage.SyncStorageImageUseCase](globalInjector)
	failedIds, err := syncUseCase.SyncStorageImages(ctx, input.StartId, input.Delta)
	if err != nil {
		return Output{}, err
	}

	return Output{FailedIds: failedIds}, nil
}

func main() {
	defer globalInjector.Shutdown()
	lambda.Start(HandleRequest)
}
