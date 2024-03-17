package conversion

import "fmt"

const n = 100

func (in *Input) Sample() *Input {
	out := &Input{}
	out.CustomFields = in.CustomFields
	out.Approaches = in.Approaches
	out.Campaigns = in.Campaigns
	out.Funds = in.Funds
	out.AttributedUserId = in.AttributedUserId
	out.CallerUserId = in.CallerUserId
	out.OrganizationAccountRecordType = in.OrganizationAccountRecordType
	out.HouseholdAccountRecordType = in.HouseholdAccountRecordType

	requiredRefs := map[string]bool{}

	addRef := func(ref *string) {
		if ref == nil {
			return
		}
		requiredRefs[*ref] = true
	}

	for _, r := range in.Relationships {
		if len(out.Relationships) < n {
			out.Relationships = append(out.Relationships, r)
			addRef(r.Account1Ref)
			addRef(r.Account2Ref)
		}
	}

	var rgs, notes, cs, gs, sc, sd int
	for _, je := range in.JournalEntries {
		if je.Disbursement != nil || je.SegmentedDonation != nil || je.Payment != nil || je.Pledge != nil {
			requiredRefs[je.Ref()] = true
			requiredRefs[je.AccountRef()] = true
			continue
		}
		if je.RecurringGift != nil && rgs < n {
			requiredRefs[je.Ref()] = true
			requiredRefs[je.AccountRef()] = true
			addRef(je.RecurringGift.RecurringGiftScheduleRef)
			if je.RecurringGift.SoftCredit != nil {
				addRef(je.RecurringGift.SoftCredit.Ref)
			}
			rgs++
		}
		if je.Note != nil && notes < n {
			requiredRefs[je.Ref()] = true
			requiredRefs[je.AccountRef()] = true
			notes++
		}
		if je.Contact != nil && cs < n {
			requiredRefs[je.Ref()] = true
			requiredRefs[je.AccountRef()] = true
			cs++
		}
		if je.Gift != nil && gs < n {
			requiredRefs[je.Ref()] = true
			requiredRefs[je.AccountRef()] = true
			if je.Gift.SoftCredit != nil {
				addRef(je.Gift.SoftCredit.Ref)
			}
			gs++
		}
		if je.SoftCredit != nil && sc < n {
			requiredRefs[je.Ref()] = true
			requiredRefs[je.AccountRef()] = true
			addRef(je.SoftCredit.AccountRef)
			addRef(je.SoftCredit.HardCreditRef)
			addRef(je.SoftCredit.HardCreditAccountRef)
			sc++
		}
		if je.SegmentedDonation != nil && sd < n {
			requiredRefs[je.Ref()] = true
			requiredRefs[je.AccountRef()] = true
			if je.Gift.SoftCredit != nil {
				addRef(je.Gift.SoftCredit.Ref)
			}
			sd++
		}
	}

	for _, je := range in.JournalEntries {
		if !requiredRefs[je.Ref()] {
			continue
		}
		if je.RecurringGiftSchedule != nil {
			addRef(je.RecurringGiftSchedule.AccountRef)
		}
		if je.SoftCredit != nil {
			addRef(je.SoftCredit.HardCreditAccountRef)
			addRef(je.SoftCredit.HardCreditRef)
			addRef(je.SoftCredit.AccountRef)
			addRef(je.SoftCredit.Ref)
		}
	}

	for _, je := range in.JournalEntries {
		if requiredRefs[je.Ref()] {
			out.JournalEntries = append(out.JournalEntries, je)
		}
	}

	seen := map[string]bool{}
	for _, a := range in.Accounts {
		if seen[*a.Ref] {
			panic(fmt.Errorf("duplicate account with ref: %s", *a.Ref))
		}
		seen[*a.Ref] = true
		if requiredRefs[*a.Ref] {
			out.Accounts = append(out.Accounts, a)
		}
	}

	fmt.Printf(`

SAMPLED
Accounts %d => %d	
JournalEntries %d => %d
Relationships %d => %d
(all others fully included)

`, len(in.Accounts), len(out.Accounts), len(in.JournalEntries), len(out.JournalEntries), len(in.Relationships), len(out.Relationships))

	return out
}
