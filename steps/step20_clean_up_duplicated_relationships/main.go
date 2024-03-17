package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Silicon-Ally/etap2sf/salesforce/clients/enterprise/utils"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
	os.Exit(0)
}

func run() error {
	client, err := utils.NewSandboxClient()
	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}
	err = client.DeleteRelationshipsNotCreatedThroughETap()
	if err != nil {
		return fmt.Errorf("deleting relationships not created through etap: %w", err)
	}
	fmt.Printf("Done cleaning up duplicated relationships. This concludes the migration. Well done!\n")
	return nil
}
