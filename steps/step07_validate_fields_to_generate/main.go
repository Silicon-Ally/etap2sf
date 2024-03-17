package main

import (
	"log"
	"os"

	"github.com/Silicon-Ally/etap2sf/conv/validate_fields_to_generate"
)

func main() {
	if err := validate_fields_to_generate.Run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
