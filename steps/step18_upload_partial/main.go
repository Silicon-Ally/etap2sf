package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Silicon-Ally/etap2sf/salesforce/upload/upload_data_to_salesforce"
)

func main() {
	if err := upload_data_to_salesforce.Run(true); err != nil {
		log.Fatal(err)
	}
	fmt.Print("\nNote: This is a partial upload. Please proceed to the next step to upload the full data set.\n\n")
	os.Exit(0)
}
