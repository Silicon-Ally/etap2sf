package generate_converters

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Silicon-Ally/etap2sf/conv/conversionsettings"
	"github.com/Silicon-Ally/etap2sf/utils"
)

func Run() error {
	for eot, sots := range conversionsettings.ObjectTypeMap {
		for _, sot := range sots {
			standard, custom, err := eTapToSalesforceFieldMappings(eot, sot)
			if err != nil {
				return fmt.Errorf("failed to construct mappings: %v", err)
			}
			code, err := transformationGoCode(eot, sot, standard, custom)
			if err != nil {
				return fmt.Errorf("failed to generate code for ETAP %s to SF %s: %v", eot, sot, err)
			}
			code, err = handleImports(code)
			if err != nil {
				return fmt.Errorf("failed to handle imports for ETAP %s to SF %s: %v", eot, sot, err)
			}
			filePath := filepath.Join(utils.ProjectRoot(), "conv", "conversion", strings.ToLower(fmt.Sprintf("generated_%s_to_%s.go", eot, sot)))
			if err := os.WriteFile(filePath, []byte(code), 0777); err != nil {
				return fmt.Errorf("failed to write code: %v", err)
			}
		}
	}
	fmt.Printf("Done generating automatic converters. However, you now must add any custom logic you want for converting in any non-standard way. You can find each of these in the `conv/conversion/manual.go` file. Once you think that is ready, you can proceed to the next step.\n")
	return nil
}
