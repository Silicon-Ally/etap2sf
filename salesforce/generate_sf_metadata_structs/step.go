package generate_sf_metadata_structs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	sfgenericclient "github.com/Silicon-Ally/etap2sf/salesforce/clients/generic/utils"
	"github.com/Silicon-Ally/etap2sf/utils"
)

var packageName = "sfmetadata"

func Run() error {
	client, err := sfgenericclient.NewGenericClient()
	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}

	if err := client.DownloadMetadataWSDL(filepath.Join(utils.ProjectRoot(), "salesforce", "generated", "sfmetadata", "sfmetadata.wsdl")); err != nil {
		return fmt.Errorf("downloading metadata wsdl: %w", err)
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

	replaceOne := map[string]string{}
	replaceAll := map[string]string{
		"\tScale int32": "\tScale *int32",
	}
	p := filepath.Join(destinationFolder, "sfmetadata", "generated-sfmetadata.go")
	eData, err := os.ReadFile(p)
	if err != nil {
		return fmt.Errorf("error reading sfmetadata generated file: %w", err)
	}
	e := string(eData)
	for repl, with := range replaceOne {
		e = strings.Replace(e, repl, with, 1)
	}
	for repl, with := range replaceAll {
		e = strings.Replace(e, repl, with, -1)
	}
	if err := os.WriteFile(p, []byte(e), 0644); err != nil {
		return fmt.Errorf("error writing sfmetadata generated file: %w", err)
	}

	return nil
}
