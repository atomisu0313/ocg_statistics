package main

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"atomisu.com/ocg-statics/cardApiCall/internal/usecase/master"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
)

type mockMasterUseCase struct {
	insertCardInfoFunc               func(ctx context.Context, neuronCardID int64) (int64, error)
	insertCardInfoListFunc           func(ctx context.Context, startId int64, delta int64) ([]int64, error)
	insertCardInfoListByOcgCardIdsFunc func(ctx context.Context, ocgCardIds []int64) ([]int64, error)
}

func (m *mockMasterUseCase) InsertCardInfo(ctx context.Context, neuronCardID int64) (int64, error) {
	if m.insertCardInfoFunc != nil {
		return m.insertCardInfoFunc(ctx, neuronCardID)
	}
	return 0, nil
}

func (m *mockMasterUseCase) InsertCardInfoList(ctx context.Context, startId int64, delta int64) ([]int64, error) {
	if m.insertCardInfoListFunc != nil {
		return m.insertCardInfoListFunc(ctx, startId, delta)
	}
	return nil, nil
}

func (m *mockMasterUseCase) InsertCardInfoListByOcgCardIds(ctx context.Context, ocgCardIds []int64) ([]int64, error) {
	if m.insertCardInfoListByOcgCardIdsFunc != nil {
		return m.insertCardInfoListByOcgCardIdsFunc(ctx, ocgCardIds)
	}
	return nil, nil
}

func TestHandleRequest_InvalidJSON(t *testing.T) {
	event := json.RawMessage(`{ invalid json }`)
	output, err := HandleRequest(context.Background(), event)
	assert.Error(t, err)
	assert.Empty(t, output.FailedIds)
}

func TestHandleRequest_EmptyOcgCardIds(t *testing.T) {
	testCases := []struct {
		name  string
		event string
	}{
		{
			name:  "empty array",
			event: `{"ocgCardIds": []}`,
		},
		{
			name:  "omitted field",
			event: `{}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := HandleRequest(context.Background(), json.RawMessage(tc.event))
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "ocgCardIds must not be empty")
			assert.Empty(t, output.FailedIds)
		})
	}
}

func TestHandleRequest_Success(t *testing.T) {
	mockUC := &mockMasterUseCase{
		insertCardInfoListByOcgCardIdsFunc: func(ctx context.Context, ocgCardIds []int64) ([]int64, error) {
			assert.Equal(t, []int64{10000, 10001}, ocgCardIds)
			return []int64{10001}, nil
		},
	}

	do.OverrideValue[master.MasterUseCase](globalInjector, mockUC)

	event := json.RawMessage(`{"ocgCardIds": [10000, 10001]}`)
	output, err := HandleRequest(context.Background(), event)
	assert.NoError(t, err)
	assert.Equal(t, []int64{10001}, output.FailedIds)
}

func TestHandleRequest_UseCaseError(t *testing.T) {
	mockUC := &mockMasterUseCase{
		insertCardInfoListByOcgCardIdsFunc: func(ctx context.Context, ocgCardIds []int64) ([]int64, error) {
			return nil, errors.New("db connection failure")
		},
	}

	do.OverrideValue[master.MasterUseCase](globalInjector, mockUC)

	event := json.RawMessage(`{"ocgCardIds": [10000]}`)
	output, err := HandleRequest(context.Background(), event)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db connection failure")
	assert.Empty(t, output.FailedIds)
}
