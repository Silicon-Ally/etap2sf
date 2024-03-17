package main

import (
	"log"
	"os"

	"github.com/Silicon-Ally/etap2sf/salesforce/generate_sf_enterprise_structs"
)

func main() {
	if err := generate_sf_enterprise_structs.Run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
