package main

import (
	"fmt"
	"log"
	"os"

	etap "github.com/Silicon-Ally/etap2sf/etap/client"
	sf "github.com/Silicon-Ally/etap2sf/salesforce/clients/generic/utils"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("There was an issue logging into one of the clients - please check your configuration, the instructions are in the README.md in this package.")
		log.Fatal(err)
	}
	os.Exit(0)
}

func run() error {
	_, err := etap.WithClient(func(c *etap.Client) ([]byte, error) {
		_, err := c.GetAllApproaches()
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		fmt.Print(etapInstructions)
		return fmt.Errorf("trouble logging into eTapestry: %w", err)
	}

	_, err = sf.NewGenericClient()
	if err != nil {
		fmt.Print(salesforceInstructions)
		return fmt.Errorf("trouble logging into Salesforce : %w", err)
	}
	fmt.Printf("Success - able to log into both clients! You can proceed to the next step\n")
	return nil
}

const etapInstructions = `
We were unable to log into eTapestry. Please follow the instructions below to set up your eTapestry API key.

If issues persist, please file bug.

In order to connect to your eTapestry data, you'll need to generate (or use) an API key.

Follow the instructions here to generate that API key:

https://webfiles-sc1.blackbaud.com/support/howto/coveo/etapestry/etapapi.html

Then, once you have it, replace the value of the /secrets/etapestry-api-key.txt file with it.

Then, rerun this step.

Details of the error:
`

const salesforceInstructions = `
We were unable to log into salesforce. 

In order to connect to your Salesforce data, you'll need to generate a security token, and enter your credentials into the right file.

Follow these instructions:

1 - Create a salesforce Sandbox (if you already have one, or if you are doing a final deployment to production, skip this step).
    (Sandboxes are temporary copies of your configuration that allow you to experiment without worry. You can delete them without harming your data, and can create new ones at any time.)

2 - Visit the sandbox and sign in (note you may have to have .sandboxname after your email)

3 - Clicking your user in the upper right hand corner, then going to settings, create a new security token. Save this for next step (you'll have to do this once per sandbox you create.)

4 - Enter your username, password, security token, and login URL in secrets/salesforce-sandbox-connection-config.txt

5 - Rerun this step until you succeed. If you have any issues, please file a bug.

`
