package exportfiles

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type JEType string

const (
	JETypeNote                  JEType = "Note"
	JETypeContact               JEType = "Contact"
	JETypeGift                  JEType = "Gift"
	JETypeRecurringGift         JEType = "RecurringGift"
	JETypeRecurringGiftSchedule JEType = "RecurringGiftSchedule"
	JETypePledge                JEType = "Pledge"
	JETypePayment               JEType = "Payment"
)

type Request struct {
	AttachmentRef   string
	JournalEntryRef string
	TargetUser      string
	FileName        string
	JEType          JEType
}

type Authn struct {
	JSessionID        string
	UserDataSessionID string
	AuthSVCToken      string
	SecurityToken     string
	MyEntityRoleRef   string
}

func (authn *Authn) CallJEPage(request *Request) error {
	var (
		u      string
		values url.Values
	)
	switch request.JEType {
	case JETypeContact:
		u = "https://bos.etapestry.com/prod/editJournalContact.do"
		values = url.Values{"journalEntryRef": {request.JournalEntryRef}}
	case JETypeNote:
		u = "https://bos.etapestry.com/prod/editJournalNote.do"
		values = url.Values{"journalEntryRef": {request.JournalEntryRef}}
	case JETypeRecurringGift, JETypeRecurringGiftSchedule, JETypeGift:
		u = "https://bos.etapestry.com/prod/editJournalTransaction.do"
		values = url.Values{
			"entityRoleRef":  {authn.MyEntityRoleRef},
			"transactionRef": {request.JournalEntryRef},
		}
	default:
		return fmt.Errorf("unsupported journal entry type: %s", request.JEType)
	}

	req, err := http.NewRequest(http.MethodGet, u+"?"+values.Encode(), nil)
	if err != nil {
		return fmt.Errorf("failed to init request: %w", err)
	}

	req.Header.Set("Authority", "bos.etapestry.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cookie", fmt.Sprintf("preferredLocale=en; JSESSIONID=%s; userDataSessionID=%s; AuthSvcToken=%s", authn.JSessionID, authn.UserDataSessionID, authn.AuthSVCToken))
	req.Header.Set("Referer", fmt.Sprintf("https://bos.etapestry.com/prod/entityRoleHome.do?entityRoleRef=%s", request.TargetUser))
	req.Header.Set("Sec-Ch-Ua", `Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "Linux")
	req.Header.Set("Sec-Fetch-Dest", "iframe")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to issue request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		dat, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-200 response code %d for JE page request: %q", resp.StatusCode, string(dat))
	}

	return nil
}

func (authn *Authn) DownloadToFile(request *Request) error {
	contactMethodRef := "000.0.00000000"

	url := "https://bos.etapestry.com/prod/downloadSessionAttachment.do"

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	var writeErr error
	writeField := func(name, val string) {
		if writeErr != nil {
			return
		}
		if err := writer.WriteField(name, val); err != nil {
			writeErr = fmt.Errorf("failed to write field %q, %q: %w", name, val, err)
		}
	}

	writeField("saveAsTemplate", "")
	writeField("name", "NA")
	writeField("primaryEmailAddress", "etap2sf")
	writeField("securityToken", authn.SecurityToken)
	writeField("attendeeRef", request.TargetUser)
	writeField("saveAndEntityRoleType", "user")
	writeField("destinationAfterSave", "persona")
	writeField("entityRoleRef", request.TargetUser)
	writeField("attendeeRef", request.TargetUser)
	writeField("journalEntryRef", request.JournalEntryRef)
	writeField("date", "9/17/2023")
	writeField("subject", "Example for Attachments")
	writeField("contactMethodRef", contactMethodRef)
	writeField("notes", "")
	writeField("attachmentFileId", request.AttachmentRef)
	writeField("attachmentIndividualLimit", "15728640")
	writeField("attachmentCollectiveLimit", "15728640")
	writeField("attachmentFileHostRef", "")

	part, err := writer.CreateFormFile("attachmentUploadFile", "")
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write([]byte{}); err != nil {
		return fmt.Errorf("failed to write form part: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close form writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return fmt.Errorf("failed to format download request: %w", err)
	}
	req.Header.Set("Authority", "bos.etapestry.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Cookie", fmt.Sprintf("preferredLocale=en; JSESSIONID=%s; userDataSessionID=%s; AuthSvcToken=%s", authn.JSessionID, authn.UserDataSessionID, authn.AuthSVCToken))
	req.Header.Set("Origin", "https://bos.etapestry.com")
	req.Header.Set("Referer", "https://bos.etapestry.com/prod/editJournalContact.do?journalEntryRef="+request.JournalEntryRef)
	req.Header.Set("Sec-Ch-Ua", `"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Fetch-Dest", "iframe")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		dat, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-200 response code %d for download request: %q", resp.StatusCode, string(dat))
	}

	f, err := os.Create(request.FileName)
	if err != nil {
		return fmt.Errorf("failed to create file for download: %w", err)
	}
	defer f.Close()

	// Download to both a file + an in-memory buffer, so we can quickly look for error info.
	var buf bytes.Buffer
	w := io.MultiWriter(&buf, f)
	if _, err := io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("failed to download to memory: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close download file: %w", err)
	}

	if bytes.Contains(buf.Bytes(), []byte("eTapestry Application Error")) {
		return fmt.Errorf("response was an eTapestryApplicationError\n\tsee the eTap error at %s", request.FileName)
	}

	return nil
}

func (a *Authn) Download(r *Request) error {
	errFile := "curl-error-output.txt"
	if err := os.MkdirAll(filepath.Dir(r.FileName), 0777); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", r.FileName, err)
	}

	if err := a.CallJEPage(r); err != nil {
		if err2 := os.WriteFile(errFile, []byte(err.Error()), 0777); err2 != nil {
			return fmt.Errorf("page call finished with error: %w, and writing error file also failed: %v", err, err2)
		}
		return fmt.Errorf("page call finished with error: %w", err)
	}

	time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)

	if err := a.DownloadToFile(r); err != nil {
		if err2 := os.WriteFile(errFile, []byte(err.Error()), 0777); err2 != nil {
			return fmt.Errorf("page call finished with error: %w, and writing error file also failed: %v", err, err2)
		}
		return fmt.Errorf("failed to download file: %w", err)
	}

	return nil
}
