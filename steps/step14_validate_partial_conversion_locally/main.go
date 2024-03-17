package main

import (
	"log"
	"os"

	"github.com/Silicon-Ally/etap2sf/conv/conversion/validate_conversion_locally"
)

func main() {
	if err := validate_conversion_locally.Run(true); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
