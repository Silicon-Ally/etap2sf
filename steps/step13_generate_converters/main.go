package main

import (
	"log"
	"os"

	"github.com/Silicon-Ally/etap2sf/conv/generate_converters"
)

func main() {
	if err := generate_converters.Run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
