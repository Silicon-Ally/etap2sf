package main

import (
	"log"
	"os"

	"github.com/Silicon-Ally/etap2sf/salesforce/create_sf_layouts"
)

func main() {
	if err := create_sf_layouts.Run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
