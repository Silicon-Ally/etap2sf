package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
	os.Exit(0)
}

func run() error {
	return fmt.Errorf(`MANUAL STEP NEEDED:

Your LAST MANUAL STEP!!! YAY!
This one matters a LOT THOUGH, so please do the entirety of this step very carefully.

Read the full post, and only move on to the next step when you're on a Salesforce page showing that you have this permission.

The real upload WILL FAIL if this is not done correctly.

https://ongkrab.medium.com/salesforce-step-assign-value-createdbyid-fields-2f007469f5b4

	`)
}
