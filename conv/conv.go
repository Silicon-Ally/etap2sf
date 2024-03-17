package conv

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/salesforce"
)

type ETapField interface {
	Name() string
	GoCodeAssignToVar(varname string) (string, error)
}

type SalesforceField interface {
	Name() string
	GoCodeAssignmentFromVar(varname string, sfot salesforce.ObjectType) (string, error)
}

type FieldMapping struct {
	In  ETapField
	Out SalesforceField
}

func (fm *FieldMapping) GoCode(varname string, sfot salesforce.ObjectType) (string, error) {
	i, err := fm.In.GoCodeAssignToVar(varname)
	if err != nil {
		return "", fmt.Errorf("writing input code: %w", err)
	}
	o, err := fm.Out.GoCodeAssignmentFromVar(varname, sfot)
	if err != nil {
		return "", fmt.Errorf("writing output code: %w", err)
	}
	return fmt.Sprintf("\t%s\n\t%s\n", i, o), nil
}
