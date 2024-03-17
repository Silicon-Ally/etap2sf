package generate_sf_enterprise_structs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	sfgenericclient "github.com/Silicon-Ally/etap2sf/salesforce/clients/generic/utils"
	"github.com/Silicon-Ally/etap2sf/utils"
)

const packageName = "sfenterprise"

func Run() error {
	client, err := sfgenericclient.NewGenericClient()
	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}

	if err := client.DownloadEnterpriseWSDL(filepath.Join(utils.ProjectRoot(), "salesforce", "generated", "sfenterprise", "sfenterprise.wsdl")); err != nil {
		return fmt.Errorf("downloading enterprise wsdl: %w", err)
	}

	// Note this goes before we generate the code from the WSDL because we modify the WSDL
	// In this step.
	if err := generateEnums(); err != nil {
		return fmt.Errorf("generating enums: %w", err)
	}

	destinationFolder := filepath.Join(utils.ProjectRoot(), "salesforce", "generated")
	wsdl := filepath.Join(destinationFolder, packageName, packageName+".wsdl")
	c := fmt.Sprintf("gowsdl -p %s -o generated-%s.go -i -d %s %s", packageName, packageName, destinationFolder, wsdl)
	cmd := exec.Command("bash", "-c", c)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fileName, _ := utils.WriteBytesToTempFile(output, "gowsdl-output-"+packageName)
		fmt.Printf("wrote gowsdl output to %s\n", fileName)
		return fmt.Errorf("command finished with error: %w", err)
	}

	replaceOne := map[string]string{
		// This is never defined anywhere - I don't get it, but it behaves like a string!
		"*QName": "string",
		// These are defined multiple times in the schema.
		"type Location ":               "type Location1 ",
		"type DescribeGlobalTheme ":    "type DescribeGlobalTheme1 ",
		"type DescribeApprovalLayout ": "type DescribeApprovalLayout1 ",
		"type DescribeLayout ":         "type DescribeLayout1 ",
		// Prevents self-package name collision.
		`"github.com/hooklift/gowsdl/soap"`: `gowsdlsoap "github.com/hooklift/gowsdl/soap"`,
		// Prevents collisions when XML Serializing objects.
		"XMLName xml.Name `xml:\"urn:enterprise.soap.sforce.com address\"`":                                                                                             "",
		"XMLName xml.Name `xml:\"urn:sobject.enterprise.soap.sforce.com npsp__Level__c\"`":                                                                              "",
		"XMLName xml.Name `xml:\"urn:sobject.enterprise.soap.sforce.com npsp__Batch__c\"`":                                                                              "",
		"XMLName xml.Name `xml:\"urn:sobject.enterprise.soap.sforce.com npo02__Household__c\"`":                                                                         "",
		"XMLName xml.Name `xml:\"urn:sobject.enterprise.soap.sforce.com npsp__Address__c\"`":                                                                            "",
		"XMLName xml.Name `xml:\"urn:sobject.enterprise.soap.sforce.com sObject\"`":                                                                                     "",
		"Npe4__ReciprocalRelationship__r *Npe4__Relationship__c `xml:\"npe4__ReciprocalRelationship__r,omitempty\" json:\"npe4__ReciprocalRelationship__r,omitempty\"`": "",
		"XMLName xml.Name `xml:\"urn:sobject.enterprise.soap.sforce.com npe03__Recurring_Donation__c\"`":                                                                "",
	}
	replaceAll := map[string]string{
		// Completes the self-package name collision fix.
		"soap.XSDDate": "gowsdlsoap.XSDDate",
		"soap.XSDTime": "gowsdlsoap.XSDTime",
		"soap.Client":  "gowsdlsoap.Client",
		// SObject Deserialization fails but we really need the ID. This is the only field we care about
		// so we just remove the nested element of the struct and take it's relevant fields to the top level.
		"\t*SObject\n": "\tId *ID `xml:\"Id,omitempty\" json:\"Id,omitempty\"`\n",
	}
	p := filepath.Join(destinationFolder, "sfenterprise", "generated-sfenterprise.go")
	eData, err := os.ReadFile(p)
	if err != nil {
		return fmt.Errorf("error reading sfenterprise generated file: %w", err)
	}
	e := string(eData)
	for repl, with := range replaceOne {
		e = strings.Replace(e, repl, with, 1)
	}
	for repl, with := range replaceAll {
		e = strings.Replace(e, repl, with, -1)
	}
	if err := os.WriteFile(p, []byte(e), 0644); err != nil {
		return fmt.Errorf("error writing sfenterprise generated file: %w", err)
	}

	fmt.Printf("You've successfully created the %s package! You can now see it in your code editor, though it may be large. You may now proceed to the next step\n", packageName)

	return nil
}
