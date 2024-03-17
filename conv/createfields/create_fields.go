package createfields

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/conv/conversionsettings"
	"github.com/Silicon-Ally/etap2sf/conv/etapfields"
	"github.com/Silicon-Ally/etap2sf/etap/inference/customfields"
	"github.com/Silicon-Ally/etap2sf/etap/inference/standardfields"
	"github.com/Silicon-Ally/etap2sf/salesforce"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfmetadata"
)

func GetSalesforceFieldsDerivedFromETapestry(sot salesforce.ObjectType) ([]*sfmetadata.CustomField, []error) {
	externalFieldKey, err := sot.SalesforceObjectExternalFieldKey()
	if err != nil {
		return nil, []error{fmt.Errorf("getting salesforce object external field key for %s: %w", sot, err)}
	}
	customFields, err := customfields.GetCustomFields()
	if err != nil {
		return nil, []error{fmt.Errorf("getting custom fields: %w", err)}
	}
	groupedByET := customFields.GroupedByETapObject()
	fields := []*sfmetadata.CustomField{}
	errors := []error{}
	for et, sts := range conversionsettings.ObjectTypeMap {
		cfs := groupedByET[et]
		found := false
		for _, st := range sts {
			if st != sot {
				continue
			}
			found = true
			for _, cf := range cfs {
				f, err := etapfields.ETapCustomFieldToSFCustomField(cf)
				if err != nil {
					errors = append(errors, fmt.Errorf("failed to convert custom field: %w", err))
				} else {
					fields = append(fields, f)
				}
			}
		}
		if !found {
			continue
		}
		if et.IsString() {
			sffm, err := standardfields.StandardFieldForString(et)
			if err != nil {
				errors = append(errors, fmt.Errorf("getting standard field for string %q: %w", et, err))
				continue
			}
			f, err := etapfields.ETapStandardFieldToSFCustomField("", sffm)
			if err != nil {
				errors = append(errors, fmt.Errorf("converting standard field for %q: %w", et, err))
				continue
			}
			fields = append(fields, f)
		} else {
			eto, err := et.Struct()
			if err != nil {
				errors = append(errors, fmt.Errorf("getting struct for %q: %w", et, err))
				continue
			}
			sfs, err := standardfields.StandardFieldsFromStruct(et, eto)
			if err != nil {
				errors = append(errors, fmt.Errorf("getting standard fields from struct %q: %w", et, err))
				continue
			}
			for _, sf := range sfs {
				f, err := etapfields.ETapStandardFieldToSFCustomField(eto, sf)
				if err != nil {
					errors = append(errors, fmt.Errorf("converting standard field for %q: %w", et, err))
					continue
				}
				if f == nil {
					continue
				}
				fields = append(fields, f)
			}
		}
	}
	if len(errors) > 0 {
		return nil, errors
	}
	if externalFieldKey == salesforce.MultiObjectExternalFieldKey {
		fields = append(fields, &sfmetadata.CustomField{
			Metadata: &sfmetadata.Metadata{
				FullName: salesforce.MultiObjectExternalFieldKey,
			},
			Label:       "Etap: Multi-Object External Field Key",
			Description: "This field is used to have a single external field key for objects that don't have a 1:1 mapping into salesforce",
			Type_:       ptr(sfmetadata.FieldTypeText),
			Length:      255,
		})
	}
	fields = append(fields, &sfmetadata.CustomField{
		Metadata: &sfmetadata.Metadata{
			FullName: "etap_MigrationExplanation__c",
		},
		Label:       "Etap: Migration Explanation",
		Description: "This field used to have a human-readable explanation of what generated this object in its current state",
		Type_:       ptr(sfmetadata.FieldTypeText),
		Length:      255,
	})
	fields = append(fields, &sfmetadata.CustomField{
		Metadata: &sfmetadata.Metadata{
			FullName: "etap_MigrationTime__c",
		},
		Label:       "Etap: Migration Time",
		Description: "This field describes the time that this entity was migrated. Useful for debugging.",
		Type_:       ptr(sfmetadata.FieldTypeDateTime),
	})
	if sot == salesforce.ObjectType_Task || sot == salesforce.ObjectType_ContentVersion {
		fields = append(fields, &sfmetadata.CustomField{
			Metadata: &sfmetadata.Metadata{
				FullName: "etap_AdditionalContextForRecord__c",
			},
			Label:                   "Etap: Additional Context",
			Description:             "This field links to another object of type 'AdditionalContext' which stores values that are truncated in this entity",
			Type_:                   ptr(sfmetadata.FieldTypeLookup),
			ReferenceTo:             "etap_AdditionalContext__c",
			RelationshipLabel:       fmt.Sprintf("%s Additional Context", sot),
			RelationshipName:        fmt.Sprintf("%s_Additional_Context", sot),
			DeleteConstraint:        ptr(sfmetadata.DeleteConstraintSetNull),
			WriteRequiresMasterRead: false,
		})
	}

	return fields, nil
}

func ptr[T any](t T) *T {
	return &t
}
