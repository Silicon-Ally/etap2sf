package upload_data_to_salesforce

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/conv/conversion"
	"github.com/Silicon-Ally/etap2sf/salesforce/upload"
)

func Run(partial bool) error {
	input, err := conversion.GetInput()
	if err != nil {
		return fmt.Errorf("failed to get input: %v", err)
	}
	if partial {
		input = input.Sample()
	}
	output, err := input.Convert()
	if err != nil {
		return fmt.Errorf("failed to convert to output: %v", err)
	}
	doShuffles := !partial
	uploader, err := upload.GetOrCreateUploader(doShuffles)
	if err != nil {
		return fmt.Errorf("getting or creating uploader: %w", err)
	}
	if err := uploader.Upload(output); err != nil {
		return fmt.Errorf("uploading: %w", err)
	}
	return nil
}
