package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Silicon-Ally/etap2sf/salesforce/generate_sf_metadata_structs"
)

func main() {
	if err := generate_sf_metadata_structs.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Done generating Salesforce metadata structs. Proceed to next step.\n")
	os.Exit(0)
}
