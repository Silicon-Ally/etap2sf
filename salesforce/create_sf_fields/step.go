package create_sf_fields

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/Silicon-Ally/etap2sf/conv/validate_fields_to_generate"
	"github.com/Silicon-Ally/etap2sf/etap/data"
	"github.com/Silicon-Ally/etap2sf/salesforce/clients/metadata/utils"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfmetadata"
	"golang.org/x/exp/maps"
)

func Run() error {
	tcs, err := validate_fields_to_generate.GetValidatedFieldsToGenerate()
	if err != nil {
		return fmt.Errorf("getting validated fields to generate: %w", err)
	}
	errs, err := readErrors()
	if err != nil {
		return fmt.Errorf("reading errors: %w", err)
	}
	defer func() {
		if err := errs.Write(); err != nil {
			log.Printf("failed to write errors to disk: %v", err)
		}
	}()

	sort.Slice(tcs, func(i, j int) bool {
		return errs.sortIdx(tcs[i].ID()) < errs.sortIdx(tcs[j].ID())
	})

	finished := 0
	needed := 0
	for _, tc := range tcs {
		if isErr, ok := errs.Ids[tc.ID()]; isErr || !ok {
			needed++
		} else {
			finished++
		}
	}

	client, err := utils.NewMetadataSandboxClient()
	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}

	relationships, err := data.GetRelationships()
	if err != nil {
		return fmt.Errorf("getting relationships: %w", err)
	}
	rts := map[string]bool{}
	for _, relationship := range relationships {
		rts[*relationship.Type.Role1] = true
		rts[*relationship.Type.Role2] = true
	}
	if err := client.AddRelationshipTypesToPicklist(maps.Keys(rts)); err != nil {
		return fmt.Errorf("adding relationship types to picklist: %w", err)
	}

	errs.Errors = []string{}
	for i, t := range tcs {
		fmt.Printf("%d + (%d / %d) %s %s\n", finished, i, needed, t.ObjectType, t.CustomField.FullName)
		if isErr, ok := errs.Ids[t.ID()]; ok && !isErr {
			fmt.Printf("  skipping - already inserted successfully\n")
			continue
		}
		err = client.UpsertCustomField(t.ObjectType, t.CustomField)
		if err != nil {
			s := fmt.Sprintf("upserting metadata for %s %s: %v", t.ObjectType, t.CustomField.FullName, err)
			errs.Errors = append(errs.Errors, s)
			errs.Ids[t.ID()] = true
			return errors.New(s)
		} else {
			errs.Ids[t.ID()] = false
		}
	}

	// The following are omitted from the list of profiles because I didn't know what they
	// were used for/if they ought have access. If you want a different set of profiles to have access, change this list.
	//
	// "B2BMA Integration User",
	// "Chatter External User",
	// "Chatter Free User",
	// "Chatter Moderator User",
	// "ContractManager",
	// "Guest",
	// "Guest License User",
	// "Identity User",
	// "Minimum Access - Salesforce",
	// "Read Only",
	// "SolutionManager",
	// "Standard",
	// "StandardAul",
	profiles := []string{
		"Admin",
		"Executive Management",
		"Fundraising and Development",
		"MarketingProfile",
		"Office Staff",
		"Salesforce API Only System Integrations",
	}
	for _, p := range profiles {
		pf := &sfmetadata.Profile{
			Metadata: &sfmetadata.Metadata{
				FullName: p,
			},
		}
		for _, t := range tcs {
			sfn, err := t.ObjectType.SalesforceNameForFieldCreation()
			if err != nil {
				return fmt.Errorf("getting salesforce name for %s: %w", t.ObjectType, err)
			}
			pf.FieldPermissions = append(pf.FieldPermissions, &sfmetadata.ProfileFieldLevelSecurity{
				Field:    fmt.Sprintf("%s.%s", sfn, t.CustomField.FullName),
				Readable: true,
				Editable: true,
			})
		}
		if err := client.UpdateProfile(pf); err != nil {
			return fmt.Errorf("updating profile: %w", err)
		}
	}

	fmt.Printf("Done - fields have been created. You may proceed to the next step.\n")
	return nil
}
