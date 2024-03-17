package exportfiles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Silicon-Ally/etap2sf/etap/data"
	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/etap/generated/overrides"
	"github.com/Silicon-Ally/etap2sf/utils"
)

var manifestPath = filepath.Join(utils.ProjectRoot(), "etap", "attachments", "download-manifest.json")
var AttachmentsFolder = filepath.Join(utils.ProjectRoot(), "etap", "attachments", "raw")

type ManifestEntry struct {
	Attachment   generated.Attachment
	JournalEntry overrides.JournalEntry
	DoneSuccess  bool
	DoneError    bool
	OutputFile   string
}

func (me *ManifestEntry) FilePath() string {
	return fmt.Sprintf("%s/%s/%s", AttachmentsFolder, *me.Attachment.Ref, *me.Attachment.Filename)
}

func GetPathToAttachment(a *generated.Attachment) (string, error) {
	if a.Ref == nil {
		return "", fmt.Errorf("attachment ref is nil")
	}
	if a.Filename == nil {
		return "", fmt.Errorf("attachment filename is nil")
	}
	return fmt.Sprintf("%s/%s/%s", AttachmentsFolder, *a.Ref, *sanitizeFilename(a.Filename)), nil
}

func (m *ManifestEntry) ValidateWithLocalFile() error {
	if _, err := os.Stat(m.FilePath()); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to stat output file: %w", err)
	}
	b, err := os.ReadFile(m.FilePath())
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	if len(b) == 0 || bytes.Contains(b, []byte("eTapestry Application Error")) {
		m.DoneError = true
		m.DoneSuccess = false
	} else {
		m.DoneError = false
		m.DoneSuccess = true
	}
	return nil
}

func (me *ManifestEntry) Process(authn *Authn) error {
	t, err := journalEntryToType(&me.JournalEntry)
	if err != nil {
		return fmt.Errorf("failed to get type: %w", err)
	}
	request := &Request{
		JournalEntryRef: me.JournalEntry.Ref(),
		TargetUser:      me.JournalEntry.AccountRef(),
		AttachmentRef:   *me.Attachment.Ref,
		FileName:        me.FilePath(),
		JEType:          t,
	}
	if err := authn.Download(request); err != nil {
		me.DoneError = true
		return fmt.Errorf("failed to download attachment: %w", err)
	}
	return me.ValidateWithLocalFile()
}

type Manifest struct {
	Entries []*ManifestEntry
	Success int
}

func (m *Manifest) Load() error {
	fmt.Println("Loading manifest...")
	bytes, err := os.ReadFile(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Creating manifest...")
			jes, err := data.GetJournalEntries()
			if err != nil {
				return fmt.Errorf("failed to get journal entries: %w", err)
			}
			utils.Shuffle(jes)
			m.Entries = []*ManifestEntry{}
			for _, je := range jes {
				for _, a := range je.Attachments() {
					me := &ManifestEntry{
						JournalEntry: *je,
						Attachment:   *a,
						DoneSuccess:  false,
						DoneError:    false,
					}
					me.Attachment.Filename = sanitizeFilename(me.Attachment.Filename)
					m.Entries = append(m.Entries, me)
				}
			}
			m.Success = 0
			for _, e := range m.Entries {
				if err := e.ValidateWithLocalFile(); err != nil {
					return fmt.Errorf("failed to validate manifest entry: %w", err)
				}
				if e.DoneSuccess {
					m.Success++
				}
			}
			return m.Save()
		}
		return fmt.Errorf("failed to read all-journal-entries.json: %w", err)
	}
	if err := json.Unmarshal(bytes, m); err != nil {
		return fmt.Errorf("failed to unmarshal manifest: %w", err)
	}
	m.Success = 0
	for _, e := range m.Entries {
		if err := e.ValidateWithLocalFile(); err != nil {
			return fmt.Errorf("failed to validate manifest entry: %w", err)
		}
		if e.DoneSuccess {
			m.Success++
		}
	}
	return nil
}

func (m *Manifest) Save() error {
	fmt.Println("Saving manifest...")
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal manifest: %w", err)
	}
	if err := os.WriteFile(manifestPath, b, 0777); err != nil {
		return fmt.Errorf("unable to write manifest: %w", err)
	}
	return nil
}

func sanitizeFilename(s *string) *string {
	if s == nil {
		return nil
	}
	var builder strings.Builder
	for _, runeValue := range *s {
		if runeValue != '$' {
			builder.WriteRune(runeValue)
		}
	}
	ss := builder.String()
	return &ss
}

func (m *Manifest) Process(authn *Authn) (bool, error) {
	for i, e := range m.Entries {
		if !e.DoneSuccess {
			fmt.Printf("%d/%d (%d OK) Processing %s\n", i, len(m.Entries), m.Success, *e.Attachment.Ref)
			if err := e.Process(authn); err != nil {
				return true, fmt.Errorf("failed to process manifest entry: %w", err)
			}
			if e.DoneSuccess {
				m.Success++
				e.DoneError = false
			}
			if e.DoneError {
				return true, fmt.Errorf("failed to process manifest entry: %s", *e.Attachment.Ref)
			}
			return true, nil
		}
	}
	return false, nil
}

func journalEntryToType(je *overrides.JournalEntry) (JEType, error) {
	if je.Note != nil {
		return JETypeNote, nil
	}
	if je.Gift != nil {
		return JETypeGift, nil
	}
	if je.Contact != nil {
		return JETypeContact, nil
	}
	if je.Pledge != nil {
		return JETypePledge, nil
	}
	if je.Payment != nil {
		return JETypePayment, nil
	}
	if je.RecurringGift != nil {
		return JETypeRecurringGift, nil
	}
	if je.RecurringGiftSchedule != nil {
		return JETypeRecurringGiftSchedule, nil
	}
	return "", fmt.Errorf("unknown journal entry type: %+v", je)
}
