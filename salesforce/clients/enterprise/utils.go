package client

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfenterprise"
	"github.com/hooklift/gowsdl/soap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func StructToFieldsMap(data interface{}) (map[string]interface{}, error) {
	toTitle := cases.Title(language.English)
	result := make(map[string]interface{})
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("data is not a struct or pointer to struct")
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		structField := t.Field(i)

		// Ignore unexported fields and XML Name Fields
		if structField.PkgPath != "" || fieldValue.Type().String() == "xml.Name" {
			continue
		}
		if !isEmptyValue(fieldValue) {
			key := structField.Name // default to using the Go struct field name
			// Check for an "xml" tag, and if present, use it as the key
			if xmlTag := structField.Tag.Get("xml"); xmlTag != "" {
				key = strings.Split(xmlTag, ",")[0] // just in case there are additional tag options
			}
			if fieldValue.Kind() == reflect.Ptr {
				fieldValue = fieldValue.Elem()
			}
			if fieldValue.Kind() == reflect.Struct {
				// Compound fields are bespoke, and require bespoke solutions, unfortunately.
				if fieldValue.Type().String() == "sfenterprise.Address" {
					fieldPrefix := strings.TrimSuffix(key, "Address")
					nestedMap, err := StructToFieldsMap(fieldValue.Interface())
					if err != nil {
						return nil, err
					}
					for k, v := range nestedMap {
						result[fieldPrefix+toTitle.String(k)] = v
					}
				} else if fieldValue.Type().String() == "soap.XSDDateTime" {
					xsddt, ok := fieldValue.Interface().(soap.XSDDateTime)
					if !ok {
						return nil, fmt.Errorf("expected soap.XSDDateTime, got %T", fieldValue.Interface())
					}
					goTime := xsddt.ToGoTime()
					if !goTime.IsZero() {
						dateTimeLayout := time.RFC3339Nano
						if goTime.Nanosecond() == 0 {
							dateTimeLayout = time.RFC3339
						}
						result[key] = goTime.Format(dateTimeLayout)
					}
				} else if fieldValue.Type().String() == "soap.XSDDate" {
					xsddt, ok := fieldValue.Interface().(soap.XSDDate)
					if !ok {
						return nil, fmt.Errorf("expected soap.XSDDate, got %T", fieldValue.Interface())
					}
					goTime := xsddt.ToGoTime()
					if !goTime.IsZero() {
						dateTimeLayout := time.RFC3339Nano
						if goTime.Nanosecond() == 0 {
							dateTimeLayout = time.RFC3339
						}
						result[key] = goTime.Format(dateTimeLayout)
					}
				} else if fieldValue.Type().String() == "sfenterprise.RecordType" {
					rt, ok := fieldValue.Interface().(sfenterprise.RecordType)
					if !ok {
						return nil, fmt.Errorf("expected sfenterprise.RecordType, got %T", fieldValue.Interface())
					}
					if rt.Id == nil {
						return nil, fmt.Errorf("expected sfenterprise.RecordType.Id to be populated, got nil")
					}
					id := string(*rt.Id)
					if id == "" {
						return nil, fmt.Errorf("expected sfenterprise.RecordType.Id to be populated, got ID=%q Name=%q SObjectType=%q RecordType=%+v", *rt.Id, *rt.Name, *rt.SobjectType, rt)
					}
					result["RecordTypeId"] = id
				} else {
					return nil, fmt.Errorf("field type %s doesn't have a clear xml serialization story", fieldValue.Type())
				}
			} else {
				if hasUnderlyingStringType(fieldValue.Type()) {
					result[key] = fieldValue.String()
				} else if st := getAsSimpleType(fieldValue); st != "" {
					result[key] = st
				} else if fieldValue.Type() == reflect.TypeOf([]uint8{}) {
					data := fieldValue.Interface().([]byte)
					if len(data) == 0 {
						return nil, fmt.Errorf("failed to encode %s as base64", key)
					}
					result[key] = string(data)
				} else {
					return nil, fmt.Errorf("field type %s doesn't have a clear xml serialization story", fieldValue.Type())
				}
			}
		}
	}
	return result, nil
}

func getAsSimpleType(v reflect.Value) string {
	switch v.Type().Kind() {
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", v.Float())
	}
	return ""
}

func hasUnderlyingStringType(t reflect.Type) bool {
	return t.Kind() == reflect.String
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
