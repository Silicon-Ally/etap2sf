package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Silicon-Ally/etap2sf/salesforce"
	"github.com/Silicon-Ally/etap2sf/salesforce/clients/metadata/utils"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
	os.Exit(0)
}

func run() error {
	client, err := utils.NewMetadataSandboxClient()
	if err != nil {
		return fmt.Errorf("getting client: %w", err)
	}
	for _, sot := range salesforce.ObjectTypes {
		if sot.IsCustomToMigration() {
			if err := client.CreateCustomObject(sot); err != nil {
				return fmt.Errorf("creating object: %w", err)
			}
		}
	}
	fmt.Print("Done creating new Salesforce objects. Proceed to next step.\n")
	return nil
}
