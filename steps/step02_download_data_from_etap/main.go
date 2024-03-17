package main

import (
	"log"
	"os"

	"github.com/Silicon-Ally/etap2sf/etap/data/download_all_data_from_etap"
)

func main() {
	if err := download_all_data_from_etap.Run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
