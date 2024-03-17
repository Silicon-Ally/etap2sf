package etapfields

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Silicon-Ally/etap2sf/conv/conversionsettings"
	"github.com/Silicon-Ally/etap2sf/etap/inference/customfields"
	"github.com/Silicon-Ally/etap2sf/salesforce"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfmetadata"
	"github.com/Silicon-Ally/etap2sf/utils"
)

type CustomField struct {
	Delegate            *customfields.CustomField
	DefinedValuesFnName string
}

func (f *CustomField) Name() string {
	return *f.Delegate.GeneratedDefinition.Name
}

func (f *CustomField) GoCodeAssignToVar(varname string) (string, error) {
	return fmt.Sprintf(`%s := GetDefinedFieldValues(%s(in), "%s")`, varname, f.DefinedValuesFnName, *f.Delegate.GeneratedDefinition.Name), nil
}

func (f *CustomField) GoCodeAssignmentFromVar(varname string, sfot salesforce.ObjectType) (string, error) {
	cf, err := ETapCustomFieldToSFCustomField(f.Delegate)
	if err != nil {
		return "", fmt.Errorf("converting custom field to salesforce custom field: %w", err)
	}
	cat := utils.AlphanumericOnly(*f.Delegate.GeneratedDefinition.Category)
	if cat == "" && *f.Delegate.GeneratedDefinition.System {
		cat = "EtapSystem"
	}
	if sub, ok := conversionsettings.CategoryNameSubstitutions[cat]; ok {
		cat = sub
	}
	n := utils.AlphanumericOnly(*f.Delegate.GeneratedDefinition.Name)
	ln := "etap_" + cat + "_" + n + "__c"
	if sub, ok := conversionsettings.FieldNameSubstitutions[ln]; ok {
		ln = sub
	}
	n = strings.Split(ln, "_")[2]

	fieldName := fmt.Sprintf("Etap_%s_%s__c", cat, n)
	typeName := fmt.Sprintf("%s_etap%s%s_", sfot, cat, n)
	switch *cf.Type_ {
	case sfmetadata.FieldTypePicklist:

		return fmt.Sprintf(`if len(%[1]s) > 1 {
			fileName, err := utils.WriteValueToTempJSONFile(%[1]s, "multiple-values-for-%[2]s")
			if err != nil {
				return nil, fmt.Errorf("dumping multiple values failed: %[3]s", err)
			}
			return nil, fmt.Errorf("multiple (%[4]s) values for %[2]s: results dumped to '%[6]s'", len(%[1]s), fileName)
		} else if len(%[1]s) == 1 {
			parsed, err := sfenterprise.Parse_%[5]s(*%[1]s[0])
			if err != nil {
				return nil, fmt.Errorf("parse error for Parse_%[5]s: %[3]s", err)
			}
			out.%[2]s = ptr(parsed)
		}`, varname, fieldName, "%w", "%d", typeName, "%s"), nil
	case sfmetadata.FieldTypeMultiselectPicklist:
		// I straight up couldn't tell you why this one only gets converted for some object entities.
		// My guess is that the WSDL refuses to generate HUGE enums. IDK though. It gets generated for AdditionalContext, so _shrug_?
		if fieldName == "Etap_EtapSystem_LogEntryType__c" && sfot != salesforce.ObjectType_AdditionalContext && sfot != salesforce.ObjectType_Task {
			return fmt.Sprintf("out.%s = ptr(JoinWithSemicolons(%s))", fieldName, varname), nil
		}
		return fmt.Sprintf(`%[1]sParsed := []sfenterprise.%[2]s{}
	for _, value := range %[1]s {
		parsed, err := sfenterprise.Parse_%[2]s(*value)
		if err != nil {
			return nil, fmt.Errorf("parse error for parse_%[2]s: %[3]s", err)
		}
		%[1]sParsed = append(%[1]sParsed, parsed)
	}
	out.%[4]s = JoinEnumsWithSemicolons(%[1]sParsed)`, varname, typeName, "%w", fieldName), nil
		// return fmt.Sprintf(`out.%s = JoinWithSemicolons(%s)`, fieldName, varname), nil
	case sfmetadata.FieldTypeText, sfmetadata.FieldTypeTextArea, sfmetadata.FieldTypeLongTextArea,
		sfmetadata.FieldTypeEmail, sfmetadata.FieldTypePhone:
		return fmt.Sprintf(`if len(%[1]s) >= 1 {
			s := JoinWithSemicolons(%[1]s)
			out.%[2]s = ptr(s)
		}`, varname, fieldName), nil
	case sfmetadata.FieldTypeDate:
		return fmt.Sprintf(`if len(%[1]s) > 1 {
		if err := dumpToTemporaryFile("multiple-values-for-%[2]s.txt", %[1]s); err != nil {
			return nil, fmt.Errorf("dumping multiple values failed: %[5]s", err)
		}	
		return nil, fmt.Errorf("multiple values for %[2]s: %[3]s", %[1]s)
	} else if len(%[1]s) == 1 {
		t, err := AttemptToParseDate(string(*%[1]s[0]))
		if err != nil {
			return nil, fmt.Errorf("at field '%[2]s': parsing time %[4]s failed: %[5]s", *%[1]s[0], err)
		}
		if t != nil && !t.IsZero() {
			out.%[2]s = ptr(soap.CreateXsdDate(*t, true))
		}
	}`, varname, fieldName, "%+v", "%s", "%w"), nil
	case sfmetadata.FieldTypeNumber, sfmetadata.FieldTypeCheckbox, sfmetadata.FieldTypeCurrency:
		parseFn := "strconv.ParseFloat"
		kind := "float64"
		secondArg := ", 64"
		trimmer := ""
		if *cf.Type_ == sfmetadata.FieldTypeCheckbox {
			parseFn = "strconv.ParseBool"
			kind = "bool"
			secondArg = ""
		} else if cf.Scale != nil && *cf.Scale == 0 {
			parseFn = "strconv.Atoi"
			kind = "int"
			secondArg = ""
		}
		if *cf.Type_ == sfmetadata.FieldTypeCurrency {
			trimmer = `
			onlyValue = strings.TrimPrefix(onlyValue, "$")
			onlyValue = strings.ReplaceAll(onlyValue, ",", "")
			`
		}
		return fmt.Sprintf(`if len(%[1]s) > 1 {
		if err := dumpToTemporaryFile("multiple-values-for-%[2]s.txt", %[1]s); err != nil {
			return nil, fmt.Errorf("dumping multiple values failed: %[5]s", err)
		}	
		return nil, fmt.Errorf("multiple values for %[2]s: %[3]s", %[1]s)
	} else if len(%[1]s) == 1 {
		onlyValue := *%[1]s[0]%[9]s
		if i, err := %[6]s(onlyValue%[8]s); err != nil {
			return nil, fmt.Errorf("converting %[2]s (%[4]s) to %[7]s: %[5]s", onlyValue, err)
		} else {
			out.%[2]s = ptr(float64(i))
		}
	}`, varname, fieldName, "%+v", "%s", "%w", parseFn, kind, secondArg, trimmer), nil
	case sfmetadata.FieldTypeDateTime:
		return "", fmt.Errorf("not implemented custom_etap_fields: %s", *cf.Type_)
	}
	return "", fmt.Errorf("unknown type %s", *cf.Type_)
}

func ETapCustomFieldToSFCustomField(cf *customfields.CustomField) (*sfmetadata.CustomField, error) {
	if len(cf.Values.InvalidValues) > 0 {
		return nil, fmt.Errorf("invalid values for %q: %v", *cf.GeneratedDefinition.Name, cf.Values.InvalidValues)
	}
	if len(cf.Values.AllValues) == 0 {
		return nil, nil
	}
	name, label, err := CreateSalesforceFieldNameAndLabel(*cf.GeneratedDefinition.Category, *cf.GeneratedDefinition.Name, *cf.GeneratedDefinition.System)
	if err != nil {
		return nil, fmt.Errorf("creating salesforce field name and label: %w", err)
	}
	result := &sfmetadata.CustomField{
		Metadata: &sfmetadata.Metadata{
			FullName: name,
		},
		Description: *cf.GeneratedDefinition.Desc,
		Label:       label,
	}

	if cf.GeneratedDefinition.Values != nil && len(cf.GeneratedDefinition.Values.Items) > 0 {
		if cf.MaxCardinality() > 1 || *cf.GeneratedDefinition.DisplayType == 2 /* MultiSelect */ {
			result.Type_ = ptr(sfmetadata.FieldTypeMultiselectPicklist)
			mc := int32(cf.MaxCardinality())
			if mc < 3 {
				mc = 3
			}
			if mc > 10 {
				mc = 10
			}
			result.VisibleLines = mc
		} else {
			result.Type_ = ptr(sfmetadata.FieldTypePicklist)
		}
		vs := &sfmetadata.ValueSet{
			ValueSetDefinition: &sfmetadata.ValueSetValuesDefinition{
				Value: []*sfmetadata.CustomValue{},
			},
		}
		for _, i := range cf.GeneratedDefinition.Values.Items {
			desc := strings.TrimSpace(*i.Value + " " + *i.Desc)
			vs.ValueSetDefinition.Value = append(vs.ValueSetDefinition.Value, &sfmetadata.CustomValue{
				Description: desc,
				IsActive:    !*i.Disabled,
				Metadata: &sfmetadata.Metadata{
					FullName: *i.Value,
				},
			})
		}
		result.ValueSet = vs
	} else {
		switch *cf.GeneratedDefinition.DataType {
		case 0:
			maxLen := 0
			for s := range cf.Values.AllValues {
				if len(s) > maxLen {
					maxLen = len(s)
				}
			}
			if isProbablyAPhoneNumberField(cf) {
				result.Type_ = ptr(sfmetadata.FieldTypePhone)
			} else if isProbablyAnEmailField(cf) {
				result.Type_ = ptr(sfmetadata.FieldTypeEmail)
			} else if maxLen <= 255 {
				if *cf.GeneratedDefinition.DisplayType == 0 {
					result.Type_ = ptr(sfmetadata.FieldTypeText)
					result.Length = 255
				} else if *cf.GeneratedDefinition.DisplayType == 3 {
					result.Type_ = ptr(sfmetadata.FieldTypeTextArea)
				} else {
					gdData, err := json.Marshal(cf.GeneratedDefinition)
					if err != nil {
						return nil, fmt.Errorf("failed to marshal generated definition: %v", err)
					}
					return nil, fmt.Errorf("unknown display type %d: %s", *cf.GeneratedDefinition.DisplayType, string(gdData))
				}
			} else if maxLen <= 32768 {
				result.Type_ = ptr(sfmetadata.FieldTypeLongTextArea)
				result.Length = 32768
				result.VisibleLines = 4
			} else {
				result.Type_ = ptr(sfmetadata.FieldTypeLongTextArea)
				result.Length = int32(float64(maxLen) * 1.25)
				result.VisibleLines = 4
			}
		case 1:
			result.Type_ = ptr(sfmetadata.FieldTypeDate)
		case 3:
			result.Type_ = ptr(sfmetadata.FieldTypeNumber)
			isInt := true
			for s := range cf.Values.AllValues {
				if _, err := strconv.ParseInt(s, 10, 64); err != nil {
					isInt = false
				}
			}
			if isInt {
				result.Scale = ptr(int32(0))
				result.Precision = 18
			} else {
				result.Scale = ptr(int32(5))
				result.Precision = 18
			}
		case 4:
			result.Type_ = ptr(sfmetadata.FieldTypeCurrency)
			result.Scale = ptr(int32(5))
			result.Precision = 18
		default:
			gdData, err := json.Marshal(cf.GeneratedDefinition)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal generated definition: %v", err)
			}
			return nil, fmt.Errorf("unsupported dataType %d: %s", *cf.GeneratedDefinition.DataType, string(gdData))
		}
	}
	return result, nil
}
