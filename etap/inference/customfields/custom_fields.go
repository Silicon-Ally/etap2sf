package customfields

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/etap"
	"github.com/Silicon-Ally/etap2sf/etap/generated"
)

type CustomFields struct {
	Fields map[string]*CustomField
}

type CustomField struct {
	GeneratedDefinition     *generated.DefinedField
	Values                  *CustomFieldValues
	OnEntityTypes           map[etap.ObjectType]bool
	OnEntityIDs             map[string]bool
	EntityIDsToUniqueValues map[string]map[string]bool
}

type CustomFieldValues struct {
	Kind          int
	AllValues     map[string]int
	InvalidValues map[string]int
}

func (cfs *CustomFields) GroupedByETapObject() map[etap.ObjectType][]*CustomField {
	grouped := map[etap.ObjectType][]*CustomField{}
	for _, df := range cfs.Fields {
		for et := range df.OnEntityTypes {
			dff := *df
			grouped[et] = append(grouped[et], &dff)
		}
	}
	return grouped
}

func (f *CustomField) Lookup(values []*generated.DefinedValue) *generated.DefinedValue {
	for _, v := range values {
		if *v.FieldRef == *f.GeneratedDefinition.Ref {
			return v
		}
	}
	return nil
}

func (d *CustomFields) LookupByNameOrFail(name string, values *generated.ArrayOfDefinedValue) (*generated.DefinedValue, error) {
	f, err := d.LookupByName(name, values)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, fmt.Errorf("field found, but no field value found on this instance for %q", name)
	}
	return f, nil
}

func (d *CustomFields) LookupByName(name string, values *generated.ArrayOfDefinedValue) (*generated.DefinedValue, error) {
	var f *CustomField
	for _, df := range d.Fields {
		if *df.GeneratedDefinition.Name == name {
			f = df
			break
		}
	}
	if f == nil {
		return nil, fmt.Errorf("no defined field %q", name)
	}
	for _, v := range values.Items {
		if *v.FieldRef == *f.GeneratedDefinition.Ref {
			return v, nil
		}
	}
	return nil, nil
}

func (d *CustomFields) LookupByRefOrFail(ref string, values *generated.ArrayOfDefinedValue) (*generated.DefinedValue, error) {
	var f *CustomField
	for _, df := range d.Fields {
		if *df.GeneratedDefinition.Ref == ref {
			f = df
			break
		}
	}
	if f == nil {
		return nil, fmt.Errorf("no defined field %q", ref)
	}
	for _, v := range values.Items {
		if *v.FieldRef == *f.GeneratedDefinition.Ref {
			return v, nil
		}
	}
	return nil, fmt.Errorf("no value for field %q", ref)
}

func (f *CustomField) MaxCardinality() int {
	max := 0
	for _, et := range f.EntityIDsToUniqueValues {
		if len(et) > max {
			max = len(et)
		}
	}
	return max
}
