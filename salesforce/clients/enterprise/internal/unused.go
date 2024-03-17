package internal

// These functions were useful at one earlier point in time, but are no longer useful.
// They are kept because they could be useful in the future, and each is a bit bespoke.

/*
func (c *Client) ListProfiles() error {
	res, err := c.metadataClient.ListMetadata([]*metaforce.ListMetadataQuery{{
		Type: "Profile",
	}})
	if err != nil {
		return fmt.Errorf("listing metadata: %w", err)
	}
	tmpFile, err := utils.WriteValueToTempFile(res, "list-profile-error")
	if err != nil {
		return fmt.Errorf("writing profile response: %w", err)
	}
	fmt.Printf("response at %s\n", tmpFile)
	return nil
}
*/

/*
func (c *Client) CreatePermissionSet(ps *sfmetadata.PermissionSet) error {
	res, err := c.metadataClient.UpsertMetadata([]metaforce.MetadataInterface{
		&struct {
			*sfmetadata.PermissionSet
			XSINS string `xml:"xmlns:xsi,attr"`
			XSIT  string `xml:"xsi:type,attr"`
		}{
			PermissionSet: ps,
			XSINS:         "http://www.w3.org/2001/XMLSchema-instance",
			XSIT:          "PermissionSet",
		},
	})
	if err != nil {
		return fmt.Errorf("inserting metadata: %w", err)
	}
	tmpFile, err := utils.WriteValueToTempFile(res, "create-permission-set-error")
	if err != nil {
		return fmt.Errorf("writing permission set: %w", err)
	}
	for _, result := range res.Result {
		if len(result.Errors) > 0 {
			return fmt.Errorf("errors found in response - see %s", tmpFile)
		}
	}
	if !res.Result[0].Success {
		return fmt.Errorf("failed to upsert field - see %s", tmpFile)
	}
	return nil
}
*/
