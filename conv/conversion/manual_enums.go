//go:build ignore_until_step_12

package conversion

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/etap/generated"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfenterprise"
)

func etapRelationshipTypeToSFRelationshipTypes(rt *generated.RelationshipType) (sfenterprise.Npe4Relationship_npe4Type_, sfenterprise.Npe4Relationship_npe4Type_, error) {
	s := fmt.Sprintf("%s / %s", *rt.Role1, *rt.Role2)
	switch s {
	case "Organization / Member":
		return sfenterprise.Npe4Relationship_npe4Type_Organization, sfenterprise.Npe4Relationship_npe4Type_Member, nil
	case "Sister / Brother-in-Law":
		return sfenterprise.Npe4Relationship_npe4Type_Sister, sfenterprise.Npe4Relationship_npe4Type_BrotherinLaw, nil
	case "Church / Youth Pastor":
		return sfenterprise.Npe4Relationship_npe4Type_Church, sfenterprise.Npe4Relationship_npe4Type_YouthPastor, nil
	case "Brother in law / Sister":
		return sfenterprise.Npe4Relationship_npe4Type_BrotherinLaw, sfenterprise.Npe4Relationship_npe4Type_SisterinLaw, nil
	case "Husband / Wife":
		return sfenterprise.Npe4Relationship_npe4Type_Husband, sfenterprise.Npe4Relationship_npe4Type_Wife, nil
	case "Sister / Sister-in-Law":
		return sfenterprise.Npe4Relationship_npe4Type_SisterinLaw, sfenterprise.Npe4Relationship_npe4Type_SisterinLaw, nil
	case "Cousin / Cousin":
		return sfenterprise.Npe4Relationship_npe4Type_Cousin, sfenterprise.Npe4Relationship_npe4Type_Cousin, nil
	case "Church Main Campus / Satellite Campus of Main Church":
		return sfenterprise.Npe4Relationship_npe4Type_ChurchMainCampus, sfenterprise.Npe4Relationship_npe4Type_SatelliteCampusofMainChurch, nil
	case "Partner / Partner":
		return sfenterprise.Npe4Relationship_npe4Type_Partner, sfenterprise.Npe4Relationship_npe4Type_Partner, nil
	case "Grandparent / Grandchild":
		return sfenterprise.Npe4Relationship_npe4Type_Grandparent, sfenterprise.Npe4Relationship_npe4Type_Grandchild, nil
	case "Mother / Daughter":
		return sfenterprise.Npe4Relationship_npe4Type_Mother, sfenterprise.Npe4Relationship_npe4Type_Daughter, nil
	case "Mother / Son":
		return sfenterprise.Npe4Relationship_npe4Type_Mother, sfenterprise.Npe4Relationship_npe4Type_Son, nil
	case "Father / Daughter":
		return sfenterprise.Npe4Relationship_npe4Type_Father, sfenterprise.Npe4Relationship_npe4Type_Daughter, nil
	case "Father / Son":
		return sfenterprise.Npe4Relationship_npe4Type_Father, sfenterprise.Npe4Relationship_npe4Type_Son, nil
	case "Mother-in-Law / Daughter or Son-in-Law":
		return sfenterprise.Npe4Relationship_npe4Type_MotherinLaw, sfenterprise.Npe4Relationship_npe4Type_DaughterorSoninLaw, nil
	case "Father-in-Law / Daughter or Son-in-Law":
		return sfenterprise.Npe4Relationship_npe4Type_FatherinLaw, sfenterprise.Npe4Relationship_npe4Type_DaughterorSoninLaw, nil
	case "Sister / Brother":
		return sfenterprise.Npe4Relationship_npe4Type_Sister, sfenterprise.Npe4Relationship_npe4Type_Brother, nil
	case "Sister / Sister":
		return sfenterprise.Npe4Relationship_npe4Type_Sister, sfenterprise.Npe4Relationship_npe4Type_Sister, nil
	case "Aunt or Uncle / Niece or Nephew":
		return sfenterprise.Npe4Relationship_npe4Type_AuntorUncle, sfenterprise.Npe4Relationship_npe4Type_NieceorNephew, nil
	case "Business / Owner":
		return sfenterprise.Npe4Relationship_npe4Type_Business, sfenterprise.Npe4Relationship_npe4Type_Owner, nil
	case "Employer / Employee":
		return sfenterprise.Npe4Relationship_npe4Type_Employer, sfenterprise.Npe4Relationship_npe4Type_Employee, nil
	case "Organization / Founder":
		return sfenterprise.Npe4Relationship_npe4Type_Organization, sfenterprise.Npe4Relationship_npe4Type_Founder, nil
	case "Business Partner / Business Partner":
		return sfenterprise.Npe4Relationship_npe4Type_BusinessPartner, sfenterprise.Npe4Relationship_npe4Type_BusinessPartner, nil
	case "Roommate / Roommate":
		return sfenterprise.Npe4Relationship_npe4Type_Roommate, sfenterprise.Npe4Relationship_npe4Type_Roommate, nil
	case "Member / Church":
		return sfenterprise.Npe4Relationship_npe4Type_Member, sfenterprise.Npe4Relationship_npe4Type_Church, nil
	case "Church / Lead Pastor":
		return sfenterprise.Npe4Relationship_npe4Type_Church, sfenterprise.Npe4Relationship_npe4Type_LeadPastor, nil
	case "Church / Associate Pastor":
		return sfenterprise.Npe4Relationship_npe4Type_Church, sfenterprise.Npe4Relationship_npe4Type_AssociatePastor, nil
	}
	return "", "", fmt.Errorf("unknown relationship type: %q", s)
}

func parsePreferredPronouns(s string) (sfenterprise.Contact_Pronouns_, error) {
	switch s {
	case "He/Him/His":
		return sfenterprise.Contact_Pronouns_HeHim, nil
	case "She/Her/Hers", "She/her/hers", "She/her":
		return sfenterprise.Contact_Pronouns_SheHer, nil
	case "They/Them/Theirs", "They/them":
		return sfenterprise.Contact_Pronouns_TheyThem, nil
	case "Prefer to self identify":
		return sfenterprise.Contact_Pronouns_NotListed, nil
	case "She/They":
		return sfenterprise.Contact_Pronouns_SheThey, nil
	}
	return "", fmt.Errorf("pronouns %q not supported", s)
}

func parseTitleToSalutation(title string) (sfenterprise.Contact_Salutation_, error) {
	switch title {
	case "Mr.":
		return sfenterprise.Contact_Salutation_Mr_, nil
	case "Dr.":
		return sfenterprise.Contact_Salutation_Dr_, nil
	case "Ms.":
		return sfenterprise.Contact_Salutation_Ms_, nil
	case "Mrs.":
		return sfenterprise.Contact_Salutation_Mrs_, nil
	case "?":
		return "", nil
	case "Pastor", "Rev. Dr.", "Reverend", "Miss":
		return sfenterprise.Contact_Salutation_(title), nil
	}
	return "", fmt.Errorf("unknown title: %q", title)
}

func salutationToGender(salutation sfenterprise.Contact_Salutation_) (sfenterprise.Contact_Gender_, sfenterprise.Contact_GenderIdentity_, error) {
	switch salutation {
	case sfenterprise.Contact_Salutation_Mr_:
		return sfenterprise.Contact_Gender_Male, sfenterprise.Contact_GenderIdentity_Male, nil
	case sfenterprise.Contact_Salutation_Ms_, sfenterprise.Contact_Salutation_Mrs_:
		return sfenterprise.Contact_Gender_Female, sfenterprise.Contact_GenderIdentity_Female, nil
	}
	return "", "", fmt.Errorf("can't determine gender from salutation: %q", salutation)
}

func pronounsToGender(pronouns sfenterprise.Contact_Pronouns_) (sfenterprise.Contact_Gender_, sfenterprise.Contact_GenderIdentity_, error) {
	switch pronouns {
	case sfenterprise.Contact_Pronouns_HeHim:
		return sfenterprise.Contact_Gender_Male, sfenterprise.Contact_GenderIdentity_Male, nil
	case sfenterprise.Contact_Pronouns_SheHer:
		return sfenterprise.Contact_Gender_Female, sfenterprise.Contact_GenderIdentity_Female, nil
	}
	return "", "", fmt.Errorf("pronouns %q not deteriminitive of gender", pronouns)
}

func convertOrganizationType(in sfenterprise.Account_etapAccountInformationAccountType_) (sfenterprise.Account_Type_, error) {
	switch in {
	case "Business":
		return sfenterprise.Account_Type_Corporate, nil
	case "Foundation":
		return sfenterprise.Account_Type_Foundation, nil
	case "Individual":
		return sfenterprise.Account_Type_Individual, nil
	case "Religious":
		return sfenterprise.Account_Type_ReligiousOrganization, nil
	case "Organization":
		return sfenterprise.Account_Type_Other, nil
	default:
		return "", fmt.Errorf("unknown organization type: %q", in)
	}
}

// https://app.etapestry.com/hosted/files/api3/objects/StandardPaymentSchedule.html
func convertInstallmentFrequency(in int) (*float64, *sfenterprise.Npe03RecurringDonation_npe03InstallmentPeriod_, error) {
	switch in {
	case 1:
		// Annually
		return ptr(12.0), ptr(sfenterprise.Npe03RecurringDonation_npe03InstallmentPeriod_Monthly), nil
	case 2:
		// Semi-Annual
		return ptr(6.0), ptr(sfenterprise.Npe03RecurringDonation_npe03InstallmentPeriod_Monthly), nil
	case 4:
		// Quarterly
		return ptr(3.0), ptr(sfenterprise.Npe03RecurringDonation_npe03InstallmentPeriod_Monthly), nil
	case 6:
		// Semi-Annual
		return ptr(2.0), ptr(sfenterprise.Npe03RecurringDonation_npe03InstallmentPeriod_Monthly), nil
	case 12:
		// Monthly
		return ptr(1.0), ptr(sfenterprise.Npe03RecurringDonation_npe03InstallmentPeriod_Monthly), nil
	case 24:
		// Semi-Monthly
		return ptr(1.0), ptr(sfenterprise.Npe03RecurringDonation_npe03InstallmentPeriod_1stand15th), nil
	case 26:
		// Bi-Weekly
		return ptr(1.0), ptr(sfenterprise.Npe03RecurringDonation_npe03InstallmentPeriod_Weekly), nil
	case 52:
		// Weekly
		return ptr(1.0), ptr(sfenterprise.Npe03RecurringDonation_npe03InstallmentPeriod_Weekly), nil
	case 101:
		// I think this is an error - its exactly one of these, and it only dontated once.
		return ptr(12.0), ptr(sfenterprise.Npe03RecurringDonation_npe03InstallmentPeriod_Monthly), nil
	}
	return nil, nil, fmt.Errorf("unknown installment frequency: %d", in)
}
