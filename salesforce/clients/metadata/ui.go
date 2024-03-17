package client

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"sort"
	"strings"

	"github.com/Silicon-Ally/etap2sf/conv/conversionsettings"
	"github.com/Silicon-Ally/etap2sf/conv/createfields"
	"github.com/Silicon-Ally/etap2sf/salesforce"
	"github.com/Silicon-Ally/etap2sf/salesforce/generated/sfmetadata"
	"github.com/google/uuid"
	"github.com/tzmfreedom/go-metaforce"
)

func (c *Client) CreateETapestrySectionsOnFlexiPages() error {
	type ReadResult struct {
		Records []*sfmetadata.FlexiPage `xml:"records,omitempty"`
	}

	type ReadMetadataResponse struct {
		XMLName xml.Name `xml:"http://soap.sforce.com/2006/04/metadata readMetadataResponse"`

		Result *ReadResult `xml:"result,omitempty"`
	}

	for _, sot := range salesforce.ObjectTypes {
		if sot == salesforce.ObjectType_ContentDocumentLink || sot == salesforce.ObjectType_ContentVersion {
			// These aren't even VISIBLE.
			// We store the additional info in the AdditionalContext object.
			continue
		}
		flexiPageNames, err := sot.SalesforceFlexiPageNames()
		if err != nil {
			return fmt.Errorf("getting flexipage names: %w", err)
		}
		visualforcePageContent, err := createApexPageForSalesforceObjectType(sot)
		if err != nil {
			return fmt.Errorf("creating visualforce page content: %w", err)
		}
		// Base64 encode visualforcePageContent
		visualforcePageName := fmt.Sprintf("etap_%s_migrated_data", sot)
		err = handleUpsert(c.gc.MetadataClient.UpsertMetadata([]metaforce.MetadataInterface{
			&struct {
				*sfmetadata.ApexPage
				XSINS string `xml:"xmlns:xsi,attr"`
				XSIT  string `xml:"xsi:type,attr"`
			}{
				ApexPage: &sfmetadata.ApexPage{
					Label:            fmt.Sprintf("eTapestry Migrated Data %s", sot),
					AvailableInTouch: true,
					Description:      fmt.Sprintf("eTapestry Migrated Data %s", sot),
					MetadataWithContent: &sfmetadata.MetadataWithContent{
						Metadata: &sfmetadata.Metadata{
							FullName: visualforcePageName,
						},
						Content: encodeToBase64(visualforcePageContent),
					},
				},
				XSINS: "http://www.w3.org/2001/XMLSchema-instance",
				XSIT:  "ApexPage",
			},
		}))
		if err != nil {
			return fmt.Errorf("upserting apex page: %w", err)
		}

		for _, flexiPageName := range flexiPageNames {
			resp := &ReadMetadataResponse{}
			err = c.gc.MetadataClient.ReadMetadataInto("FlexiPage", []string{flexiPageName}, resp)
			if err != nil {
				return fmt.Errorf("reading flexi page %s: %w", flexiPageName, err)
			}
			if resp.Result == nil || len(resp.Result.Records) == 0 {
				return fmt.Errorf("no flexipages found with name %q", flexiPageName)
			}
			if len(resp.Result.Records) > 1 {
				return fmt.Errorf("multiple flexipages found with name %q", flexiPageName)
			}
			if resp.Result.Records[0] == nil {
				return fmt.Errorf("no flexipage found with name %q", flexiPageName)
			}
			flexiPage := resp.Result.Records[0]
			uuid := uuid.New().String()
			facetID := "facet-" + uuid

			var main *sfmetadata.FlexiPageRegion
			var tabRegion *sfmetadata.FlexiPageRegion
			var foundETap bool
			for _, r := range flexiPage.FlexiPageRegions {
				for _, ii := range r.ItemInstances {
					if ii.ComponentInstance != nil {
						for _, p := range ii.ComponentInstance.ComponentInstanceProperties {
							switch p.Value {
							case "eTapestry", "eTapestry Migrated Data":
								foundETap = true
							case "Standard.Tab.relatedLists":
								tabRegion = r
							}
						}
					}
				}
				if r.Name == "main" {
					main = r
				}
			}
			if foundETap {
				continue
			}
			if tabRegion != nil {
				tabRegion.ItemInstances = append(tabRegion.ItemInstances, &sfmetadata.ItemInstance{
					ComponentInstance: &sfmetadata.ComponentInstance{
						ComponentInstanceProperties: []*sfmetadata.ComponentInstanceProperty{
							{
								Name:  "title",
								Value: "eTapestry",
							},
							{
								Name:  "body",
								Value: facetID,
							},
							{
								Name:  "active",
								Value: "false",
							},
						},
						ComponentName: "flexipage:tab",
						Identifier:    "etapMigrationTab",
					},
				})
				flexiPage.FlexiPageRegions = append(flexiPage.FlexiPageRegions, &sfmetadata.FlexiPageRegion{
					Type_: ptr(sfmetadata.FlexiPageRegionTypeFacet),
					Name:  facetID,
					ItemInstances: []*sfmetadata.ItemInstance{{
						ComponentInstance: &sfmetadata.ComponentInstance{
							ComponentName: "flexipage:visualforcePage",
							Identifier:    "flexipage_visualforcePage_etapTab",
							ComponentInstanceProperties: []*sfmetadata.ComponentInstanceProperty{{
								Name: "height", Value: "1000",
							}, {
								Name: "label", Value: "eTapestry Migrated Data",
							}, {
								Name: "showLabel", Value: "false",
							}, {
								Name: "pageName", Value: visualforcePageName,
							}},
						},
					}},
				})
			} else {
				main.ItemInstances = append(main.ItemInstances, &sfmetadata.ItemInstance{
					ComponentInstance: &sfmetadata.ComponentInstance{
						ComponentName: "flexipage:visualforcePage",
						Identifier:    "flexipage_visualforcePage_etapMain",
						ComponentInstanceProperties: []*sfmetadata.ComponentInstanceProperty{{
							Name: "height", Value: "1000",
						}, {
							Name: "label", Value: "eTapestry Migrated Data",
						}, {
							Name: "showLabel", Value: "false",
						}, {
							Name: "pageName", Value: visualforcePageName,
						}},
					},
				})
			}
			err = handleUpdate(c.gc.MetadataClient.UpdateMetadata([]metaforce.MetadataInterface{
				&struct {
					*sfmetadata.FlexiPage
					XSINS string `xml:"xmlns:xsi,attr"`
					XSIT  string `xml:"xsi:type,attr"`
				}{
					FlexiPage: flexiPage,
					XSINS:     "http://www.w3.org/2001/XMLSchema-instance",
					XSIT:      "FlexiPage",
				},
			}))
			if err != nil {
				if strings.Contains(err.Error(), "CANNOT_MODIFY_MANAGED_OBJECT") {
					continue
				} else {
					return fmt.Errorf("updating flexipage: %w", err)
				}
			}
		}
	}
	return nil
}

func createApexPageForSalesforceObjectType(sot salesforce.ObjectType) (string, error) {
	nColumns := 2
	if sot == salesforce.ObjectType_AdditionalContext {
		nColumns = 1
	}
	sotName, err := sot.SalesforceName()
	if err != nil {
		return "", fmt.Errorf("getting salesforce name: %w", err)
	}
	fields, errs := createfields.GetSalesforceFieldsDerivedFromETapestry(sot)
	if len(errs) > 0 {
		return "", fmt.Errorf("found %d errors, the first one: %w", len(errs), errs[0])
	}
	fieldsGroupedByLabel := map[string][]*sfmetadata.CustomField{}
	for _, field := range fields {
		splits := strings.Split(field.Label, ":")
		sectionLabel := strings.Join(splits[:len(splits)-1], ":")
		sectionLabel = strings.TrimSpace(sectionLabel)
		sectionLabel = strings.ReplaceAll(sectionLabel, "ETap:", "eTapestry:")
		if sub := conversionsettings.SectionLabelSubstitutions[sectionLabel]; sub != "" {
			sectionLabel = sub
		}
		fieldsGroupedByLabel[sectionLabel] = append(fieldsGroupedByLabel[sectionLabel], field)
	}
	sortFn := func(cfs []*sfmetadata.CustomField) func(i, j int) bool {
		return func(i, j int) bool {
			return cfs[i].Label < cfs[j].Label
		}
	}
	makeSection := func(fields []*sfmetadata.CustomField, label string) (string, string, error) {
		sort.Slice(fields, sortFn(fields))
		format := `<apex:pageBlockSection columns="%d" title="%s" %s >
				%s
			</apex:pageBlockSection>`
		inputLines := []string{}
		outputLines := []string{}
		sectionTests := []string{}
		for _, field := range fields {
			fullFieldName := fmt.Sprintf("%s.%s", sotName, field.FullName)
			var test string
			switch *field.Type_ {
			case sfmetadata.FieldTypeCheckbox:
				test = ""
			case sfmetadata.FieldTypeCurrency, sfmetadata.FieldTypeNumber:
				test = fmt.Sprintf("AND(NOT(ISNULL( %[1]s )), NOT( %[1]s == 0 ))", fullFieldName)
			case sfmetadata.FieldTypeDate, sfmetadata.FieldTypeDateTime:
				test = fmt.Sprintf("NOT(ISNULL(%s))", fullFieldName)
			case sfmetadata.FieldTypeEmail, sfmetadata.FieldTypePhone, sfmetadata.FieldTypePicklist, sfmetadata.FieldTypeText,
				sfmetadata.FieldTypeTextArea, sfmetadata.FieldTypeLookup,
				sfmetadata.FieldTypeUrl, sfmetadata.FieldTypeMultiselectPicklist, sfmetadata.FieldTypeLongTextArea:
				test = fmt.Sprintf("NOT(%s == '')", fullFieldName)
			default:
				return "", "", fmt.Errorf("metadata-ui unsupported field type %s", *field.Type_)
			}
			if test != "" {
				sectionTests = append(sectionTests, test)
				test = fmt.Sprintf(`rendered="{!%s}"`, test)
			}
			outputLines = append(outputLines, fmt.Sprintf(
				`<apex:outputField value="{!%[1]s.%[2]s}" %[3]s />`,
				sotName, field.FullName, test))
			inputLines = append(inputLines, fmt.Sprintf(
				`<apex:inputField value="{!%[1]s.%[2]s}" />`,
				sotName, field.FullName))
		}
		sectionTest := ""
		if len(sectionTests) > 0 {
			sectionTest = fmt.Sprintf(`rendered="{!%s}"`, strings.Join(sectionTests, " || "))
		}
		inputSection := fmt.Sprintf(format, 1, label, "", strings.Join(inputLines, "\n\t\t\t\t"))
		outputSection := fmt.Sprintf(format, nColumns, label, sectionTest, strings.Join(outputLines, "\n\t\t\t\t"))
		return inputSection, outputSection, nil
	}

	labels := []string{}
	for label := range fieldsGroupedByLabel {
		labels = append(labels, label)
	}
	sort.Strings(labels)
	newInputSections := []string{}
	newOutputSections := []string{}
	for _, label := range labels {
		i, o, err := makeSection(fieldsGroupedByLabel[label], label)
		if err != nil {
			return "", fmt.Errorf("making section for label %s: %w", label, err)
		}
		newInputSections = append(newInputSections, i)
		newOutputSections = append(newOutputSections, o)
	}

	return fmt.Sprintf(`<apex:page standardController="%s" extensions="EditButton">  
	<p>
		This tab includes data migrated automatically from eTapestry. Editing this data directly is not advised,
		since these fields dont have integrations with other fields in Salesforce.
	</p>
	<apex:form id="etap-migration-form">
		<apex:commandButton value="{!IF(isEditing, 'Cancel', 'Edit')}" action="{!toggleIsEditing}" rerender="etap-migration-form" />
		<apex:commandButton value="Save" action="{!save}" rendered="{!isEditing}" rerender="etap-migration-form" />
		<apex:pageBlock rendered="{!isEditing}">
    		%s
    	</apex:pageBlock>
    	<apex:pageBlock rendered="{!NOT(isEditing)}">
        <h1>Note - you are not in editing mode, empty fields are omitted.</h1>
			%s
    	</apex:pageBlock>
  	</apex:form>
</apex:page>`, sotName, strings.Join(newInputSections, "\n\t\t\t"), strings.Join(newOutputSections, "\n\t\t\t\t")), nil
}

func encodeToBase64(str string) []byte {
	data := []byte(str)
	encodedLen := base64.StdEncoding.EncodedLen(len(data))
	encodedData := make([]byte, encodedLen)
	base64.StdEncoding.Encode(encodedData, data)
	return encodedData
}
