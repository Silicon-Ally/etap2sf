package create_sf_layouts

import (
	"fmt"

	mc "github.com/Silicon-Ally/etap2sf/salesforce/clients/metadata/utils"
)

func Run() error {
	client, err := mc.NewMetadataSandboxClient()
	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}
	if err := client.CreateETapestrySectionsOnFlexiPages(); err != nil {
		return fmt.Errorf("creating eTapestry sections on flexi pages: %w", err)
	}
	fmt.Printf("Success - You've successfully created visualforce components for showing migrated eTapestry data. You may proceed to the next step.\n")
	return nil
}
