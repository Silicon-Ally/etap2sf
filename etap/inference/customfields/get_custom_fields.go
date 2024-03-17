package customfields

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Silicon-Ally/etap2sf/etap"
	"github.com/Silicon-Ally/etap2sf/etap/data"
	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/etap/generated/overrides"
	"github.com/Silicon-Ally/etap2sf/utils"
)

func GetCustomFields() (*CustomFields, error) {
	cfData, err := utils.MemoizeOperation("etap-inference-custom-fields.json", func() ([]byte, error) {
		definedFields, err := data.GetDefinedFields()
		if err != nil {
			return nil, fmt.Errorf("failed to get defined fields: %w", err)
		}
		jes, err := data.GetJournalEntries()
		if err != nil {
			return nil, fmt.Errorf("failed to get journal entries: %w", err)
		}
		accounts, err := data.GetAccounts()
		if err != nil {
			return nil, fmt.Errorf("failed to get accounts: %w", err)
		}
		relationships, err := data.GetRelationships()
		if err != nil {
			return nil, fmt.Errorf("failed to get relationships: %w", err)
		}

		dfs := newCustomFields()
		for _, df := range definedFields {
			dfs.addField(df)
		}

		for _, je := range jes {
			if err := processJE(dfs, je); err != nil {
				return nil, fmt.Errorf("processing journal entry %q: %w", je.Ref(), err)
			}
		}
		for _, a := range accounts {
			for _, dv := range a.AccountDefinedValues.Items {
				if err := dfs.addValue(dv, "Account", *a.Ref); err != nil {
					return nil, fmt.Errorf("adding account %q: %w", *a.Ref, err)
				}
			}
			for _, dv := range a.PersonaDefinedValues.Items {
				if err := dfs.addValue(dv, "Persona", *a.Ref); err != nil {
					return nil, fmt.Errorf("adding persona %q: %w", *a.Ref, err)
				}
			}
		}
		for _, r := range relationships {
			for _, dv := range r.DefinedValues.Items {
				if err := dfs.addValue(dv, "Relationship", *r.Ref); err != nil {
					return nil, fmt.Errorf("adding relationship %q: %w", *r.Ref, err)
				}
			}
		}

		result, err := json.MarshalIndent(dfs, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal defined fields: %w", err)
		}
		return result, nil
	})
	if err != nil {
		return nil, fmt.Errorf("memoizing custom fields: %w", err)
	}
	result := &CustomFields{}
	if err := json.Unmarshal(cfData, result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal custom fields: %w", err)
	}
	return result, nil
}

func processJE(dfs *CustomFields, j *overrides.JournalEntry) error {
	if j.Note != nil {
		for _, dv := range j.Note.DefinedValues.Items {
			if err := dfs.addValue(dv, etap.ObjectType_Note, *j.Note.Ref); err != nil {
				return fmt.Errorf("adding value to defined fields: %w", err)
			}
		}
	}
	if j.Contact != nil {
		for _, dv := range j.Contact.DefinedValues.Items {
			if err := dfs.addValue(dv, etap.ObjectType_Contact, *j.Contact.Ref); err != nil {
				return fmt.Errorf("adding value to defined fields: %w", err)
			}
		}
	}
	if j.Gift != nil {
		for _, dv := range j.Gift.DefinedValues.Items {
			if err := dfs.addValue(dv, etap.ObjectType_Gift, *j.Gift.Ref); err != nil {
				return fmt.Errorf("adding value to defined fields: %w", err)
			}
		}
	}
	if j.Pledge != nil {
		for _, dv := range j.Pledge.DefinedValues.Items {
			if err := dfs.addValue(dv, etap.ObjectType_Pledge, *j.Pledge.Ref); err != nil {
				return fmt.Errorf("adding value to defined fields: %w", err)
			}
		}
	}
	if j.Payment != nil {
		for _, dv := range j.Payment.DefinedValues.Items {
			if err := dfs.addValue(dv, etap.ObjectType_Payment, *j.Payment.Ref); err != nil {
				return fmt.Errorf("adding value to defined fields: %w", err)
			}
		}
	}
	if j.RecurringGiftSchedule != nil {
		for _, dv := range j.RecurringGiftSchedule.DefinedValues.Items {
			if err := dfs.addValue(dv, etap.ObjectType_RecurringGiftSchedule, *j.RecurringGiftSchedule.Ref); err != nil {
				return fmt.Errorf("adding value to defined fields: %w", err)
			}
		}
	}
	if j.RecurringGift != nil {
		for _, dv := range j.RecurringGift.DefinedValues.Items {
			if err := dfs.addValue(dv, etap.ObjectType_RecurringGift, *j.RecurringGift.Ref); err != nil {
				return fmt.Errorf("adding value to defined fields: %w", err)
			}
		}
	}
	if j.Disbursement != nil {
		for _, dv := range j.Disbursement.DefinedValues.Items {
			if err := dfs.addValue(dv, etap.ObjectType_Disbursement, *j.Disbursement.Ref); err != nil {
				return fmt.Errorf("adding value to defined fields: %w", err)
			}
		}
	}
	if j.Purchase != nil {
		for _, dv := range j.Purchase.DefinedValues.Items {
			if err := dfs.addValue(dv, etap.ObjectType_Purchase, *j.Purchase.Ref); err != nil {
				return fmt.Errorf("adding value to defined fields: %w", err)
			}
		}
	}
	if j.Invitation != nil && j.Invitation.CalendarItem != nil {
		return fmt.Errorf("invitation calendar item not supported because we don't have any in our data, but was requested here")
		/*
			for _, dv := range j.Invitation.CalendarItem.DefinedValues.Items {
				if err := dfs.addValue(dv, "InvitationCalendarItem", *j.Invitation.Ref); err != nil {
					return fmt.Errorf("adding value to defined fields: %w", err)
				}
			}
		*/
	}
	return nil
}

func (dv *CustomFieldValues) add(v *generated.DefinedValue, df *CustomField) error {
	vv := ""
	if v.Value != nil {
		vv = *v.Value
	}
	valid, err := dv.isValidValue(v.Value, df)
	if err != nil {
		return fmt.Errorf("checking if value %q is valid: %w", vv, err)
	}
	dv.AllValues[*v.Value]++
	if !valid {
		dv.InvalidValues[vv]++
		return nil
	}
	return nil
}

var dateRegex = regexp.MustCompile("^[0-9]{1,2}/[0-9]{1,2}/[0-9]{4}$")

func (dv *CustomFieldValues) isValidValue(v *string, df *CustomField) (bool, error) {
	if valid, err := dv.isValidValueByType(v); err != nil {
		return false, fmt.Errorf("checking if value %q is valid by type: %w", *v, err)
	} else if !valid {
		return false, nil
	}
	return df.isValidValueByValue(v)
}

func (df *CustomField) isValidValueByValue(v *string) (bool, error) {
	dt := *df.GeneratedDefinition.DisplayType
	if dt == 0 || dt == 3 {
		return true, nil // text or text area
	}
	if dt < 0 || dt > 3 {
		return false, fmt.Errorf("unsupported dataType %d", dt)
	}
	for _, i := range df.GeneratedDefinition.Values.Items {
		if *i.Value == *v {
			return true, nil
		}
	}
	return false, nil
}

func (dv *CustomFieldValues) isValidValueByType(v *string) (bool, error) {
	switch dv.Kind {
	case 0:
		return v != nil, nil
	case 1:
		return v != nil && dateRegex.MatchString(*v), nil
	case 2:
		return false, fmt.Errorf("unsupported dataType %d", dv.Kind)
	case 3:
		return isValidInt(v), nil
	case 4:
		if v == nil {
			return false, nil
		}
		vv := *v
		if strings.HasPrefix(vv, "$") {
			vv = vv[1:]
		} else {
			return false, nil
		}
		vv = strings.ReplaceAll(vv, ",", "")
		_, err := strconv.ParseFloat(vv, 64)
		if err == nil {
			return true, nil
		} else {
			return false, nil
		}
	default:
		return false, fmt.Errorf("unsupported dataType %d", dv.Kind)
	}
}

func isValidInt(s *string) bool {
	if s == nil {
		return false
	}
	_, err := strconv.Atoi(*s)
	return err == nil
}

func newCustomField(g *generated.DefinedField) *CustomField {
	return &CustomField{
		GeneratedDefinition: g,
		Values: &CustomFieldValues{
			Kind:          *g.DataType,
			AllValues:     map[string]int{},
			InvalidValues: map[string]int{},
		},
		OnEntityTypes:           map[etap.ObjectType]bool{},
		OnEntityIDs:             map[string]bool{},
		EntityIDsToUniqueValues: map[string]map[string]bool{},
	}
}

func (d *CustomField) add(v *generated.DefinedValue, onEntityType etap.ObjectType, onEntityId string) error {
	if !ptrEq(d.GeneratedDefinition.DataType, v.DataType) {
		return fmt.Errorf("adding value %q (%s) of type %v to field %q of type %v", *v.FieldName, *v.FieldRef, *v.DataType, *d.GeneratedDefinition.Name, *d.GeneratedDefinition.DataType)
	}
	if d.EntityIDsToUniqueValues[onEntityId] == nil {
		d.EntityIDsToUniqueValues[onEntityId] = map[string]bool{}
	}
	d.OnEntityIDs[onEntityId] = true
	d.OnEntityTypes[onEntityType] = true
	d.EntityIDsToUniqueValues[onEntityId][*v.Value] = true
	if err := d.Values.add(v, d); err != nil {
		return fmt.Errorf("adding value to fields: %w", err)
	}
	return nil
}

func ptrEq[T comparable](a, b *T) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func newCustomFields() *CustomFields {
	return &CustomFields{
		Fields: map[string]*CustomField{},
	}
}

func (d *CustomFields) addValue(v *generated.DefinedValue, onEntityType etap.ObjectType, onEntityId string) error {
	if d.Fields[*v.FieldRef] == nil {
		return fmt.Errorf("no defined field %q for value %q (%s) on %s", *v.FieldRef, *v.FieldName, *v.FieldRef, onEntityType)
	}
	return d.Fields[*v.FieldRef].add(v, onEntityType, onEntityId)
}

func (d *CustomFields) addField(f *generated.DefinedField) {
	if d.Fields[*f.Ref] == nil {
		d.Fields[*f.Ref] = newCustomField(f)
	}
}
