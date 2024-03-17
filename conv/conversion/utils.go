//go:build ignore_until_step_12

package conversion

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfenterprise"
	"github.com/Silicon-Ally/etap2sf/utils"
	"github.com/hooklift/gowsdl/soap"
)

func GetDefinedFieldValues(dvs []*generated.DefinedValue, name string) []*string {
	seen := map[string]bool{}
	result := []*string{}
	for _, dv := range dvs {
		if *dv.FieldName == name {
			v := dv.Value
			if v == nil || seen[*v] {
				continue
			}
			result = append(result, v)
			seen[*v] = true
		}
	}
	return result
}

func GetDefinedValuesForAccount(a *generated.Account) []*generated.DefinedValue {
	result := []*generated.DefinedValue{}
	if a.PersonaDefinedValues != nil {
		result = append(result, a.PersonaDefinedValues.Items...)
	}
	if a.AccountDefinedValues != nil {
		result = append(result, a.AccountDefinedValues.Items...)
	}
	return result
}

func GetDefinedValuesForNote(n *generated.Note) []*generated.DefinedValue {
	return n.DefinedValues.Items
}

func GetDefinedValuesForContact(c *generated.Contact) []*generated.DefinedValue {
	return c.DefinedValues.Items
}

func GetDefinedValuesForGift(g *generated.Gift) []*generated.DefinedValue {
	return g.DefinedValues.Items
}

func GetDefinedValuesForRecurringGift(rg *generated.RecurringGift) []*generated.DefinedValue {
	return rg.DefinedValues.Items
}

func JoinWithSemicolons(in []*string) string {
	i2 := make([]string, len(in))
	for i, s := range in {
		i2[i] = *s
	}
	return strings.Join(i2, ";")
}

func JoinEnumsWithSemicolons[T ~string](in []T) *T {
	if len(in) == 0 {
		return nil
	}
	i2 := make([]string, len(in))
	for i, s := range in {
		i2[i] = string(s)
	}
	return ptr(T(strings.Join(i2, ";")))
}

func (io *io) lookupSFAccountByRef(ref string) (*sfenterprise.Account, error) {
	if a, ok := io.out.accountsByRefs[ref]; ok && a != nil {
		return a, nil
	}
	return nil, fmt.Errorf("couldn't find account with ref %q", ref)
}

func (io *io) lookupSFContactByRef(ref string) (*sfenterprise.Contact, error) {
	if c, ok := io.out.contactsByRefs[ref]; ok && c != nil {
		return c, nil
	}
	return nil, fmt.Errorf("couldn't find contact with ref %q", ref)
}

func dumpToTemporaryFile(fileName string, obj any) error {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling %T: %w", obj, err)
	}
	if filePath, err := utils.WriteBytesToTempFile(data, fileName); err != nil {
		return fmt.Errorf("writing %q: %w", fileName, err)
	} else {
		fmt.Printf("wrote to %s\n", filePath)
	}
	return nil
}

var dateFormats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05.000Z",
	"01/02/2006",
	"1/2/2006",
	"02/2006",
	"2/2006",
}

func AttemptToParseDate(s string) (*time.Time, error) {
	for _, format := range dateFormats {
		t, err := time.Parse(format, s)
		if err == nil {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("couldn't parse %q as a date", s)
}

func AttemptToParseNilableDate[T ~string](s *T) (*soap.XSDDate, error) {
	if s != nil && *s != "" {
		t, err := AttemptToParseDate(string(*s))
		if err != nil {
			return nil, err
		}
		return ptr(soap.CreateXsdDate(*t, true)), nil
	} else {
		return nil, nil
	}
}

func AttemptToParseNilableDateTime[T ~string](s *T) (*soap.XSDDateTime, error) {
	if s != nil && *s != "" {
		t, err := AttemptToParseDate(string(*s))
		if err != nil {
			return nil, err
		}
		return ptr(soap.CreateXsdDateTime(*t, true)), nil
	} else {
		return nil, nil
	}
}

func NowXSD() *soap.XSDDateTime {
	return ptr(soap.CreateXsdDateTime(time.Now(), true))
}

func cleanEmails(s *string) []string {
	if s == nil {
		return []string{}
	}
	ss := strings.TrimSpace(*s)
	if ss == "" {
		return []string{}
	}
	// ETap used to use , for delim, and now uses ;
	ss = strings.ReplaceAll(ss, ";", " ")
	ss = strings.ReplaceAll(ss, ",", " ")
	emails := strings.Split(ss, " ")
	for i := range emails {
		emails[i] = strings.TrimSpace(emails[i])
	}
	return emails
}

func clonePtr[T any](t *T) *T {
	if t == nil {
		return nil
	}
	return ptr(*t)
}

func ptr[T any](t T) *T {
	return &t
}

func trimIfLongerThan(sp *string, limit int) *string {
	if sp == nil {
		return nil
	}
	s := strings.TrimSpace(*sp)
	if len(s) == 0 {
		return nil
	}
	if len(s) > limit {
		s = s[:limit-3] + "..."
	}
	return &s
}

func errIfLongerThan(sp *string, limit int) (*string, error) {
	if sp == nil {
		return nil, nil
	}
	s := strings.TrimSpace(*sp)
	if len(s) == 0 {
		return nil, nil
	}
	if len(s) > limit {
		return nil, fmt.Errorf("string is too long (%d > %d): %q", len(s), limit, s)
	}
	return &s, nil
}

func required[T any](t *T, err error) (*T, error) {
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, fmt.Errorf("required value is nil")
	}
	return t, nil
}

type IsMissingHardCredit struct{ msg string }

func (e *IsMissingHardCredit) Error() string {
	return e.msg
}

func (e *IsMissingHardCredit) Is(target error) bool {
	_, ok := target.(*IsMissingHardCredit)
	return ok
}
