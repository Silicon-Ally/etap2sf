package validate_conversion_locally

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/conv/conversion"
	"github.com/Silicon-Ally/etap2sf/conv/generate_converters"
	"github.com/Silicon-Ally/etap2sf/salesforce/upload"
)

func Run(partial bool) error {
	if err := generate_converters.Run(); err != nil {
		return fmt.Errorf("running generate_converters: %v", err)
	}
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

	uploader, err := upload.GetUploaderForLocalValidation()
	if err != nil {
		return fmt.Errorf("failed to get uploader: %v", err)
	}
	if err := uploader.Upload(output); err != nil {
		return fmt.Errorf("failed to upload: %v", err)
	}
	fmt.Printf("Conversion succeeded, you successfully (FAKE) uploaded %d records. You can proceed to the next step.\n", len(uploader.Succeeded))
	return nil
}
