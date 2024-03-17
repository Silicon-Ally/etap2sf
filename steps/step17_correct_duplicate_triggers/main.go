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
	err = client.FixNPSPTriggers()
	if err != nil {
		return fmt.Errorf("listing triggers: %w", err)
	}
	fmt.Printf("Done fixing NPSP triggers. Proceed to next step.\n")
	return nil
}
