package etapfields

import (
	"strings"
	"unicode"

	"github.com/Silicon-Ally/etap2sf/etap/inference/customfields"
)

func isProbablyAPhoneNumberField(df *customfields.CustomField) bool {
	if *df.GeneratedDefinition.DataType != 0 {
		return false
	}
	probYes := 0.0
	probNo := 0.0
	for v, n := range df.Values.AllValues {
		if isProbablyAPhoneNumber(v) {
			probYes += float64(n)
		} else {
			probNo += float64(n)
		}
	}
	return probYes/(probYes+probNo) > 0.95
}

func isProbablyAPhoneNumber(s string) bool {
	digits := 0
	for _, r := range s {
		if unicode.IsDigit(r) {
			digits++
		}
	}
	return digits >= 8 && digits <= 10
}

func isProbablyAnEmailField(df *customfields.CustomField) bool {
	if *df.GeneratedDefinition.DataType != 0 {
		return false
	}
	probYes := 0.0
	probNo := 0.0
	for v, n := range df.Values.AllValues {
		if len(strings.Split(v, "@")) == 2 {
			probYes += float64(n)
		} else {
			probNo += float64(n)
		}
	}
	return probYes/(probYes+probNo) > 0.95
}
