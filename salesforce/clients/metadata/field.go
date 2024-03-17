package client

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/salesforce"

	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfmetadata"
	"github.com/Silicon-Ally/etap2sf/utils"
	"github.com/tzmfreedom/go-metaforce"
)

func (c *Client) UpsertCustomField(objectType salesforce.ObjectType, cf *sfmetadata.CustomField) error {
	ots, err := objectType.SalesforceNameForFieldCreation()
	if err != nil {
		return fmt.Errorf("getting salesforce name: %w", err)
	}

	cf2, err := utils.CloneJSON(cf)
	if err != nil {
		return fmt.Errorf("cloning custom field: %w", err)
	}
	cf2.FullName = fmt.Sprintf("%s.%s", ots, cf2.FullName)
	err = handleUpsert(c.gc.MetadataClient.UpsertMetadata([]metaforce.MetadataInterface{
		&struct {
			*sfmetadata.CustomField
			XSINS string `xml:"xmlns:xsi,attr"`
			XSIT  string `xml:"xsi:type,attr"`
		}{
			CustomField: cf2,
			XSINS:       "http://www.w3.org/2001/XMLSchema-instance",
			XSIT:        "CustomField",
		},
	}))
	if err != nil {
		return fmt.Errorf("inserting metadatametadata: %w", err)
	}
	return nil
}

func (c *Client) UpdateProfile(ps *sfmetadata.Profile) error {
	err := handleUpdate(c.gc.MetadataClient.UpdateMetadata([]metaforce.MetadataInterface{
		&struct {
			*sfmetadata.Profile
			XSINS string `xml:"xmlns:xsi,attr"`
			XSIT  string `xml:"xsi:type,attr"`
		}{
			Profile: ps,
			XSINS:   "http://www.w3.org/2001/XMLSchema-instance",
			XSIT:    "Profile",
		},
	}))
	if err != nil {
		return fmt.Errorf("updating metadata: %w", err)
	}
	return nil
}
