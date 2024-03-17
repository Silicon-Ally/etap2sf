package etapfields

import (
	"fmt"
	"reflect"

	"github.com/Silicon-Ally/etap2sf/etap/inference/standardfields"
	"github.com/Silicon-Ally/etap2sf/salesforce"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfmetadata"
)

type StandardField struct {
	Delegate *standardfields.StandardField
}

func (f *StandardField) Name() string {
	return f.Delegate.FullName
}

func (f *StandardField) GoCodeAssignToVar(varname string) (string, error) {
	return f.Delegate.GoCodeAssignToVar(varname)
}

func (f *StandardField) GoCodeAssignmentFromVar(varname string, sfot salesforce.ObjectType) (string, error) {
	return f.Delegate.GoCodeAssignmentFromVar(varname)
}

func ETapStandardFieldToSFCustomField(i any, ef *standardfields.StandardField) (*sfmetadata.CustomField, error) {
	var name, fType string
	if ef.InputObjectIndex == -1 {
		name = "Name"
	} else {
		val := reflect.ValueOf(i)
		if val.Kind() != reflect.Struct {
			return nil, fmt.Errorf("expected a struct, got %v", val.Kind())
		}
		f := val.Type().Field(ef.InputObjectIndex)
		name = f.Name
		fType = f.Type.String()
	}
	fullName, label, err := CreateSalesforceFieldNameAndLabel(ef.ObjectType.String(), name, false)
	if err != nil {
		return nil, fmt.Errorf("creating salesforce field name and label: %w", err)
	}
	result := &sfmetadata.CustomField{
		Metadata: &sfmetadata.Metadata{
			FullName: fullName,
		},
		Label:       label,
		Description: fmt.Sprintf("Auto-generated eTapestry %s field %s", ef.ObjectType, name),
	}

	if ef.InputObjectIndex == -1 {
		if ef.MaxLength() > 255 {
			result.Type_ = ptr(sfmetadata.FieldTypeLongTextArea)
			result.VisibleLines = 4
			result.Length = getLongTextLength(ef.MaxLength())
		} else {
			result.Type_ = ptr(sfmetadata.FieldTypeText)
			result.Length = 255
		}
		return result, nil
	}
	switch fType {
	case "*string":
		result.Type_ = ptr(sfmetadata.FieldTypeText)
		result.Length = 255
		if ef.MaxLength() > 255 {
			result.Type_ = ptr(sfmetadata.FieldTypeLongTextArea)
			result.VisibleLines = 4
			result.Length = getLongTextLength(ef.MaxLength())
		} else {
			result.Type_ = ptr(sfmetadata.FieldTypeText)
			result.Length = 255
		}
	case "*int":
		result.Type_ = ptr(sfmetadata.FieldTypeNumber)
		result.Scale = ptr(int32(0))
		result.Precision = 18
	case "*float64":
		result.Type_ = ptr(sfmetadata.FieldTypeNumber)
		result.Scale = ptr(int32(5))
		result.Precision = 18
	case "*bool":
		result.Type_ = ptr(sfmetadata.FieldTypeCheckbox)
		result.DefaultValue = "false"
		if ef.BoolUsage[true] > ef.BoolUsage[false] {
			result.DefaultValue = "true"
		}
	case "*generated.DateTime":
		result.Type_ = ptr(sfmetadata.FieldTypeDateTime)
	case "*generated.RelationshipType":
		// This is handled by a custom converter.
		return nil, nil
	case "*generated.ArrayOfDefinedValue", "*generated.ArrayOfDefinedFieldValue":
		// Skip these are handled in the definedvalues package.
		return nil, nil
	case "*generated.ArrayOfanyType",
		"*generated.ArrayOfstring",
		"*generated.ArrayOfPhone",
		"*generated.ArrayOfSocialMediaProfile",
		"*generated.ArrayOfAttachment",
		"*generated.GeneratedReceipt",
		"*generated.OrderDetail",
		"*generated.OrderInfo",
		"*generated.SoftCredit",
		"*generated.Valuable",
		"*generated.CustomPaymentSchedule",
		"*generated.StandardPaymentSchedule":
		result.Type_ = ptr(sfmetadata.FieldTypeLongTextArea)
		result.VisibleLines = 4
		result.Length = getLongTextLength(ef.MaxLength())
	default:
		return nil, fmt.Errorf("standard-to-custom unsupported type %q", fType)
	}

	return result, nil
}

func ptr[T any](t T) *T {
	return &t
}

func getLongTextLength(maxLength int) int32 {
	if maxLength <= 256 {
		return 256
	}
	var l int32 = 256
	for l < int32(maxLength) {
		l *= 2
	}
	return l
}
