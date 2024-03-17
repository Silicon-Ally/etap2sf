package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/Silicon-Ally/etap2sf/etap/data"
	"github.com/Silicon-Ally/etap2sf/salesforce/clients/metadata/utils"
)

func main() {
	msg, err := run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(msg)
	os.Exit(0)
}

func run() (string, error) {
	client, err := utils.NewMetadataSandboxClient()
	if err != nil {
		return "", fmt.Errorf("getting client: %w", err)
	}
	customTypes, err := getCustomTaskTypes()
	if err != nil {
		return "", fmt.Errorf("getting custom task types: %w", err)
	}
	baseSOAPUrl, err := url.Parse(client.GetURL())
	if err != nil {
		return "", fmt.Errorf("parsing url: %w", err)
	}
	base := baseSOAPUrl.Scheme + "://" + baseSOAPUrl.Host
	taskSubjectURL := fmt.Sprintf("%s/lightning/setup/ObjectManager/Task/FieldsAndRelationships/Subject/addPicklistValues?tid=00T&pt=7", base)
	taskTypeURL := fmt.Sprintf("%s/lightning/setup/ObjectManager/Task/FieldsAndRelationships/Type/addPicklistValues?tid=00T&pt=7", base)

	message := fmt.Sprintf(`

These tasks unfortunately can't be automated, but this set of instructions should get you there quickly.

First, make sure you adjust some NPSP settings to what you want them to be.

1. Update the Opportunity:Donation stages to be what you want them to be.
%s/lightning/setup/ObjectManager/Opportunity/FieldsAndRelationships/StageName/view

2. Decide whether you want to have recurring donations auto-create installments or not. If you want to disable them, do so via
NPSP Settings > Recurring Donations > Installment Opportunity Auto-Creation > Disable All Installments

3. In order to allow activity tracking for eTapestry Contact + Note Types, you'll need to add them to the Task.Subject and Task.Type picklists in Salesforce.

Open the following two URLs in your browser:

%s

%s

On each, enter the following list:

%s

4. That is it! You can go to the next step once this is complete.

`, base, taskSubjectURL, taskTypeURL, strings.Join(customTypes, "\n"))

	return message, nil
}

func getCustomTaskTypes() ([]string, error) {
	jes, err := data.GetJournalEntries()
	if err != nil {
		return nil, fmt.Errorf("getting jes: %w", err)
	}
	taskTypes := map[string]bool{"Note": true}
	for _, je := range jes {
		if je.Contact == nil {
			continue
		}
		taskTypes[*je.Contact.Method] = true
	}
	result := make([]string, 0, len(taskTypes))
	for tt := range taskTypes {
		result = append(result, tt)
	}
	sort.Strings(result)
	return result, nil
}
