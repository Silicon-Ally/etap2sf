package validate_fields_to_generate

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/conv/conversionsettings"
	"github.com/Silicon-Ally/etap2sf/conv/createfields"
	"github.com/Silicon-Ally/etap2sf/salesforce"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfmetadata"
	"github.com/Silicon-Ally/etap2sf/utils"
)

func Run() error {
	fields, err := GetValidatedFieldsToGenerate()
	if err != nil {
		return fmt.Errorf("generating/validating fields to generate: %w", err)
	}
	filepath, err := utils.WriteValueToTempXMLFile(fields, "validate_fields_to_generate")
	if err != nil {
		return fmt.Errorf("writing fields to temp file: %w", err)
	}
	fmt.Printf("Generated successfully - review fields that will be generated at %s.\nOnce happy with the fields that will be generated, you may proceed to the next step.", filepath)
	return nil
}

type FieldToCreate struct {
	ObjectType  salesforce.ObjectType
	CustomField *sfmetadata.CustomField
}

func (ftc *FieldToCreate) ID() string {
	return fmt.Sprintf("%s-%s", ftc.ObjectType, ftc.CustomField.FullName)
}

func GetValidatedFieldsToGenerate() ([]*FieldToCreate, error) {
	tcs := []*FieldToCreate{}
	errors := []error{}
	for _, sot := range salesforce.ObjectTypes {
		// These cannot be modified.
		if sot == salesforce.ObjectType_ContentDocumentLink {
			continue
		}
		sotETapKey, err := sot.SalesforceObjectExternalFieldKey()
		if err != nil {
			errors = append(errors, fmt.Errorf("getting salesforce object external field key for %s: %w", sot, err))
			continue
		}

		fs, errs := createfields.GetSalesforceFieldsDerivedFromETapestry(sot)
		if len(errs) > 0 {
			for _, err := range errs {
				errors = append(errors, fmt.Errorf("getting salesforce fields for %s: %w", sot, err))
			}
			continue
		}
		fs = append(fs, conversionsettings.CustomFieldsForObjectType(sot)...)
		foundKey := false
		for _, f := range fs {
			// See https://help.salesforce.com/s/articleView?id=000385982&type=1
			if sot == salesforce.ObjectType_Task && *f.Type_ == sfmetadata.FieldTypeLongTextArea {
				fieldType := sfmetadata.FieldTypeTextArea
				f.Type_ = &fieldType
				f.Length = 0
				f.VisibleLines = 0
			}
			if sotETapKey == f.FullName {
				f.ExternalId = true
				f.Unique = true
				f.CaseSensitive = false
				foundKey = true
			}
			tcs = append(tcs, &FieldToCreate{ObjectType: sot, CustomField: f})
		}

		if !foundKey {
			errors = append(errors, fmt.Errorf("failed to find key %q for %s", sotETapKey, sot))
			continue
		}
	}
	if len(errors) > 0 {
		filepath, err := utils.WriteErrorsToTempFile(errors, "validate_fields_to_generate_errors")
		if err != nil {
			return nil, fmt.Errorf("writing errors to temp file: %w", err)
		}
		return nil, fmt.Errorf("errors found, see %s", filepath)
	}
	return tcs, nil
}
