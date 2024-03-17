package standardfields

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Silicon-Ally/etap2sf/conv/conversionsettings"
	"github.com/Silicon-Ally/etap2sf/etap"
	"github.com/Silicon-Ally/etap2sf/etap/data"
)

type StandardField struct {
	ObjectType       etap.ObjectType
	FullName         string
	Type             string
	InputObjectIndex int
	StringUsage      map[string]int
	BoolUsage        map[bool]int
	EverSet          bool
}

func (f *StandardField) MaxLength() int {
	max := 0
	for s := range f.StringUsage {
		if len(s) > max {
			max = len(s)
		}
	}
	return max
}

func (f *StandardField) GoCodeAssignToVar(varname string) (string, error) {
	return fmt.Sprintf("%s := in.%s", varname, f.FullName), nil
}

func (f *StandardField) GoCodeAssignmentFromVar(varname string) (string, error) {
	outFieldName := "etap_" + f.ObjectType.String() + "_" + f.FullName + "__c"
	if sub := conversionsettings.FieldNameSubstitutions[outFieldName]; sub != "" {
		outFieldName = sub
	}
	outFieldName = strings.ToUpper(outFieldName[:1]) + outFieldName[1:]
	switch f.Type {
	case "*string":
		return fmt.Sprintf(`if %s != nil && *%s != "" {
		out.%s = %s
	}`, varname, varname, outFieldName, varname), nil
	case "*int":
		return fmt.Sprintf(`if %s != nil && *%s != 0 {
		out.%s = ptr(float64(*%s))
	}`, varname, varname, outFieldName, varname), nil
	case "*float64":
		return fmt.Sprintf(`if %s != nil && *%s != 0 {
		out.%s = %s
	}`, varname, varname, outFieldName, varname), nil
	case "*bool":
		return fmt.Sprintf(`if %s != nil {
		out.%s = %s
	}`, varname, outFieldName, varname), nil
	case "*generated.DateTime":
		return fmt.Sprintf(`if %s != nil && *%s != "" {
		t, err := time.Parse(time.RFC3339, string(*%s))
		if err != nil {
			return nil, fmt.Errorf("at field '%s': parsing time %s failed: %s", *%s, err)
		}
		if !t.IsZero() {
			out.%s = ptr(soap.CreateXsdDateTime(t, true))
		}
	}`, varname, varname, varname, f.FullName, "%q", "%w", varname, outFieldName), nil
	case "*generated.ArrayOfDefinedValue", "*generated.ArrayOfDefinedFieldValue":
		// These are handled in the definedvalues package,
		return "", fmt.Errorf("array of defined value not supported - this shouldn't be called")
	case "*generated.Valuable",
		"*generated.GeneratedReceipt",
		"*generated.OrderDetail",
		"*generated.SoftCredit",
		"*generated.OrderInfo",
		"*generated.CustomPaymentSchedule",
		"*generated.RelationshipType",
		"*generated.StandardPaymentSchedule":
		return fmt.Sprintf(`if %[1]s != nil {
			data, err := json.Marshal(%[1]s)
			if err != nil {
				return nil, fmt.Errorf("at field '%[2]s': json stringify failed: %[3]s", err)
			}
			out.%[2]s = ptr(string(data))
		}`, varname, outFieldName, "%w"), nil
	case
		"*generated.ArrayOfAttachment",
		"*generated.ArrayOfPhone",
		"*generated.ArrayOfSocialMediaProfile",
		"*generated.ArrayOfanyType":
		return fmt.Sprintf(`if %s != nil && len(%s.Items) > 0 {
		data, err := json.Marshal(%s.Items)
		if err != nil {
			return nil, fmt.Errorf("at field '%s': json stringify failed: %s", err)
		}
		out.%s = ptr(string(data))
	}`, varname, varname, varname, f.FullName, "%w", outFieldName), nil
	case "*generated.ArrayOfstring":
		return fmt.Sprintf(`if %s != nil && len(%s.Items) > 0 {
		out.%s = ptr(strings.Join(%s.Items, ";"))
	}`, varname, varname, outFieldName, varname), nil
	default:
		return "", fmt.Errorf("code-assign-from-var unsupported type %q", f.Type)
	}
}

func StandardFieldForString(o etap.ObjectType) (*StandardField, error) {
	sf := &StandardField{
		ObjectType:       o,
		FullName:         o.String() + "Name",
		Type:             "*string",
		InputObjectIndex: -1,
		StringUsage:      map[string]int{},
	}
	if err := sf.ingestAllData(); err != nil {
		return nil, fmt.Errorf("ingesting all data: %w", err)
	}
	if !sf.EverSet {
		return nil, errors.New("standard field was never set")
	}
	return sf, nil
}

func StandardFieldsFromStruct(it etap.ObjectType, i any) ([]*StandardField, error) {
	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %v", val.Kind())
	}
	results := []*StandardField{}
	vTyp := val.Type()
	for idx := 0; idx < val.NumField(); idx++ {
		r, err := standardFieldFromStructField(it, idx, vTyp.Field(idx))
		if err != nil {
			return nil, fmt.Errorf("standard field from struct field %d: %w", idx, err)
		}
		if r != nil {
			results = append(results, r)
		}
	}
	return results, nil
}

func standardFieldFromStructField(o etap.ObjectType, fieldIndex int, f reflect.StructField) (*StandardField, error) {
	fType := f.Type.String()
	fName := f.Name

	if fType == "*generated.ArrayOfDefinedValue" || fType == "*generated.ArrayOfDefinedFieldValue" || fType == "*generated.RelationshipType" {
		// These are handled in the definedvalues package,
		return nil, nil
	}
	sf := &StandardField{
		ObjectType:       o,
		FullName:         fName,
		Type:             fType,
		InputObjectIndex: fieldIndex,
		StringUsage:      map[string]int{},
		BoolUsage:        map[bool]int{},
		EverSet:          false,
	}
	if err := sf.ingestAllData(); err != nil {
		return nil, fmt.Errorf("ingesting all data: %w", err)
	}
	if !sf.EverSet {
		return nil, nil
	}
	return sf, nil
}

func (f *StandardField) ingestAllData() error {
	d, err := data.Get(f.ObjectType)
	if err != nil {
		return fmt.Errorf("getting data: %w", err)
	}
	if f.ObjectType.IsString() {
		return f.ingestAllStrings(d)
	} else {
		return f.ingestAllStructs(d)
	}
}

func (f *StandardField) ingestAllStrings(vs []any) error {
	for _, v := range vs {
		s, ok := v.(string)
		if !ok {
			return fmt.Errorf("expected a string, got %T", v)
		}
		if s == "" {
			continue
		}
		f.addToStringUsage(s)
	}
	return nil
}

func (f *StandardField) ingestAllStructs(vs []any) error {
	for _, i := range vs {
		if err := f.ingestStruct(i); err != nil {
			return fmt.Errorf("ingesting %v: %w", i, err)
		}
	}
	return nil
}

func (f *StandardField) ingestStruct(a any) error {
	val := reflect.ValueOf(a).Elem()
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, got %v", val.Kind())
	}
	field := val.Field(f.InputObjectIndex)
	switch field.Type().String() {
	case "*string":
		f.addToStringUsage(field.Elem().String())
		return nil
	case "*generated.DateTime":
		s := field.Elem().String()
		if s != "" {
			f.EverSet = true
		}
		return nil
	case "*bool":
		if field.IsNil() {
			return nil
		} else {
			f.EverSet = true
		}
		s := field.Elem().Bool()
		f.BoolUsage[s] = f.BoolUsage[s] + 1
		return nil
	case "*int":
		if field.IsNil() {
			return nil
		}
		i := field.Elem().Int()
		if i != 0 {
			f.EverSet = true
		}
		return nil
	case "*float64":
		if field.IsNil() {
			return nil
		}
		i := field.Elem().Float()
		if i != 0 {
			f.EverSet = true
		}
		return nil
	case "*generated.ArrayOfDefinedValue", "*generated.ArrayOfDefinedFieldValue", "*[]uint8":
		return nil
	case
		"*generated.ArrayOfAttachment",
		"*generated.ArrayOfPhone",
		"*generated.ArrayOfSocialMediaProfile",
		"*generated.ArrayOfanyType",
		"*generated.ArrayOfstring",
		"*generated.CustomPaymentSchedule",
		"*generated.GeneratedReceipt",
		"*generated.OrderDetail",
		"*generated.OrderInfo",
		"*generated.RelationshipType",
		"*generated.SoftCredit",
		"*generated.StandardPaymentSchedule",
		"*generated.Valuable":
		data, err := json.Marshal(field.Interface())
		if err != nil {
			return fmt.Errorf("json stringify failed: %w", err)
		}
		s := string(data)
		if len(s) > 2 { // Can be empty brackets
			f.EverSet = true
		}
		f.addToStringUsage(s)
		return nil
	}
	return fmt.Errorf("unsupported type str-value: %q", field.Type().String())
}

func (f *StandardField) addToStringUsage(s string) {
	if s != "" {
		f.EverSet = true
		f.StringUsage[s] = f.StringUsage[s] + 1
	}
}
