package etapfields

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/conv/conversionsettings"
	"github.com/Silicon-Ally/etap2sf/utils"
)

func CreateSalesforceFieldNameAndLabel(etapCategoryName, etapFieldName string, isSystemField bool) (string, string, error) {
	sanitizedCategoryName := utils.AlphanumericOnly(etapCategoryName)
	if sub, ok := conversionsettings.CategoryNameSubstitutions[sanitizedCategoryName]; ok {
		sanitizedCategoryName = sub
	}
	if sanitizedCategoryName == "" && isSystemField {
		sanitizedCategoryName = "EtapSystem"
	}
	name := "etap_" + sanitizedCategoryName + "_" + utils.AlphanumericOnly(etapFieldName) + "__c"
	name = conversionsettings.OverrideFieldNameValue(name)
	if sub, ok := conversionsettings.FieldNameSubstitutions[name]; ok {
		name = sub
	}
	if len(name) > 40 {
		return "", "", fmt.Errorf("resulting name %q is too long (%d > 40)", name, len(name))
	}
	labelCat := etapCategoryName
	if sub := conversionsettings.CategoryLabelSubstitutions[etapCategoryName]; sub != "" {
		labelCat = sub
	}
	label := "ETap: " + labelCat + ": " + etapFieldName
	label = conversionsettings.OverrideFieldLabelValue(label)
	if sub, ok := conversionsettings.FieldLabelSubstitutions[label]; ok {
		label = sub
	}
	if len(label) > 40 {
		return "", "", fmt.Errorf("label %q is too long (%d > 40) - alter the replacement scheme in conversionsettings to fix", label, len(label))
	}
	return name, label, nil
}
