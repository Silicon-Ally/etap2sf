package main

import (
	"log"
	"os"

	"github.com/Silicon-Ally/etap2sf/salesforce/create_sf_fields"
)

func main() {
	if err := create_sf_fields.Run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
