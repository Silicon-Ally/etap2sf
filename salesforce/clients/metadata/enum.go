package client

import (
	"encoding/xml"
	"fmt"

	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfmetadata"
	"github.com/Silicon-Ally/etap2sf/utils"
	"github.com/tzmfreedom/go-soapforce"
)

func (c *Client) GetAllEnumValues() (map[string]map[string][]*sfmetadata.PicklistEntry, error) {
	metadata, err := c.gc.EnterpriseClient.DescribeGlobal()
	if err != nil {
		return nil, fmt.Errorf("describing global: %w", err)
	}
	result := map[string]map[string][]*sfmetadata.PicklistEntry{}
	for i, sobj := range metadata.Sobjects {
		fmt.Printf("describing sobject %d of %d: %s\n", i+1, len(metadata.Sobjects), sobj.Name)
		object, err := c.describeSObject(sobj.Name)
		if err != nil {
			return nil, fmt.Errorf("describing sobject %s: %w", sobj.Name, err)
		}
		result[sobj.Name] = map[string][]*sfmetadata.PicklistEntry{}
		for _, field := range object.Fields {
			if len(field.PicklistValues) > 0 {
				plvs := []*sfmetadata.PicklistEntry{}
				for _, plv := range field.PicklistValues {
					converted := &sfmetadata.PicklistEntry{
						Active:       plv.Active,
						DefaultValue: plv.DefaultValue,
						Label:        plv.Label,
						Value:        plv.Value,
					}
					plvs = append(plvs, converted)
				}
				result[sobj.Name][field.Name] = plvs
			}
		}
	}
	return result, nil
}

func (c *Client) describeSObject(name string) (*soapforce.DescribeSObjectResult, error) {
	data, err := utils.MemoizeOperation("enums/salesforce-describe-sobject-"+name+".json", func() ([]byte, error) {
		object, err := c.gc.EnterpriseClient.DescribeSObject(name)
		if err != nil {
			return nil, fmt.Errorf("describing sobject %s: %w", name, err)
		}
		data, err := xml.MarshalIndent(object, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshaling object: %w", err)
		}
		return data, nil
	})
	if err != nil {
		return nil, fmt.Errorf("memoizing operation: %w", err)
	}
	result := &soapforce.DescribeSObjectResult{}
	if err := xml.Unmarshal(data, result); err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	return result, nil
}
