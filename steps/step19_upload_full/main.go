package main

import (
	"log"
	"os"

	"github.com/Silicon-Ally/etap2sf/salesforce/upload/upload_data_to_salesforce"
)

func main() {
	if err := upload_data_to_salesforce.Run(false); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
