package data

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/etap/generated"
)

var attachments []*generated.Attachment

func GetAttachments() ([]*generated.Attachment, error) {
	if attachments != nil {
		return attachments, nil
	}
	result := []*generated.Attachment{}
	jes, err := GetJournalEntries()
	if err != nil {
		return nil, fmt.Errorf("getting jes: %w", err)
	}
	for _, je := range jes {
		result = append(result, je.Attachments()...)
	}
	attachments = result
	return result, nil
}
