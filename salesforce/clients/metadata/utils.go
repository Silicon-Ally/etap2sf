package client

import (
	"fmt"

	"github.com/tzmfreedom/go-metaforce"
)

func ptr[T any](t T) *T {
	return &t
}

func handleUpdate(res *metaforce.UpdateMetadataResponse, err error) error {
	if err != nil {
		return fmt.Errorf("pure error: %w", err)
	}
	for i, result := range res.Result {
		if len(result.Errors) > 0 {
			return fmt.Errorf("failed to update - %d errors were found in result %d code[0]=%s - see debug response", len(result.Errors), i, result.Errors[0].StatusCode)
		}
		if !result.Success {
			return fmt.Errorf("failed to update - success was false in result %d - see debug response", i)
		}
	}
	return nil
}

func handleUpsert(res *metaforce.UpsertMetadataResponse, err error) error {
	if err != nil {
		return fmt.Errorf("pure error: %w", err)
	}
	for i, result := range res.Result {
		if len(result.Errors) > 0 {
			return fmt.Errorf("failed to update - %d errors were found in result %d - see debug response", len(result.Errors), i)
		}
		if !result.Success {
			return fmt.Errorf("failed to update - success was false in result %d - see debug response", i)
		}
	}
	return nil
}
