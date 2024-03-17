package client

import (
	"fmt"

	"github.com/Silicon-Ally/etap2sf/salesforce"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfmetadata"
	"github.com/tzmfreedom/go-metaforce"
)

func (c *Client) CreateCustomObject(sot salesforce.ObjectType) error {
	fullName, err := sot.SalesforceName()
	if err != nil {
		return fmt.Errorf("getting name: %w", err)
	}
	err = handleUpsert(c.gc.MetadataClient.UpsertMetadata([]metaforce.MetadataInterface{
		&struct {
			*sfmetadata.CustomObject
			XSINS string `xml:"xmlns:xsi,attr"`
			XSIT  string `xml:"xsi:type,attr"`
		}{
			CustomObject: &sfmetadata.CustomObject{
				Metadata:    &sfmetadata.Metadata{FullName: fullName},
				Label:       string(sot),
				PluralLabel: string(sot) + "s",
				NameField: &sfmetadata.CustomField{
					Metadata: &sfmetadata.Metadata{FullName: "Name"},
					Type_:    ptr(sfmetadata.FieldTypeText),
					Label:    "Name",
				},
				DeploymentStatus: ptr(sfmetadata.DeploymentStatusDeployed),
				SharingModel:     ptr(sfmetadata.SharingModelReadWrite),
			},
			XSINS: "http://www.w3.org/2001/XMLSchema-instance",
			XSIT:  "CustomObject",
		},
	}))
	if err != nil {
		return fmt.Errorf("upserting object: %w", err)
	}
	return nil
}
