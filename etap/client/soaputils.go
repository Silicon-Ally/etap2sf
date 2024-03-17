package client

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func resolveXmlHrefs(xmlStr string) (string, error) {
	type Node struct {
		XMLName  xml.Name
		Attrs    []xml.Attr `xml:",any,attr"`
		Content  string     `xml:",chardata"`
		Children []Node     `xml:",any"`
	}

	// Parse the XML
	var doc Node
	if err := xml.Unmarshal([]byte(xmlStr), &doc); err != nil {
		return "", err
	}

	// Collect all elements with ids
	idMap := make(map[string]Node)
	var collectIds func(node Node)
	collectIds = func(node Node) {
		for _, attr := range node.Attrs {
			if attr.Name.Local == "id" {
				idMap[attr.Value] = node
			}
		}
		for _, child := range node.Children {
			collectIds(child)
		}
	}
	collectIds(doc)

	// Resolve hrefs
	var resolve func(node *Node) error
	resolve = func(node *Node) error {
		for i, attr := range node.Attrs {
			if attr.Name.Local == "href" {
				ref, ok := idMap[attr.Value[1:]]
				if !ok {
					return errors.New("unresolved href: " + attr.Value)
				}
				node.Children = append(node.Children, ref.Children...)
				node.Attrs = append(node.Attrs[:i], node.Attrs[i+1:]...)
				node.Attrs = append(node.Attrs, ref.Attrs...)
			}
		}
		for i := range node.Children {
			if err := resolve(&node.Children[i]); err != nil {
				return err
			}
		}
		node.Content = strings.TrimSpace(node.Content)
		return nil
	}

	if err := resolve(&doc); err != nil {
		return "", err
	}

	// Convert back to XML
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(doc); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func resolveXmlHrefsHttpResponse(resp *http.Response) error {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}
	resolved, err := resolveXmlHrefs(string(data))
	if err != nil {
		return fmt.Errorf("resolving xml references: %w", err)
	}
	resp.Body = io.NopCloser(bytes.NewBufferString(resolved))
	return nil
}

func addRequiredSOAPEncodingStyle(req *http.Request) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	req.Body.Close()

	toFind := `<SOAP-ENV:Envelope`
	toAppend := ` SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"`

	fixedBody := bytes.Replace(body, []byte(toFind), []byte(toFind+toAppend), -1)

	req.Body = io.NopCloser(bytes.NewBuffer(fixedBody))
	req.ContentLength = int64(len(fixedBody))
	req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))

	return nil
}

func checkForFaultCode(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}
	defer resp.Body.Close()
	// Replace the response body
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	fsStart := []byte("<faultstring>")
	fsEnd := []byte("</faultstring>")
	if bytes.Contains(body, fsStart) {
		start := bytes.Index(body, fsStart) + len(fsStart)
		end := bytes.Index(body, fsEnd)
		/* If you're running into issues, uncomment the following lines to print the response

		err = os.WriteFile("error-response.xml", body, 0777)
		if err != nil {
		 		panic(err)
		}
		*/
		return fmt.Errorf("fault code found: %q uncomment above to print detailed response", body[start:end])
	}
	/* If you're running into issues, uncomment the following lines to print the response

	n := rand.Int63()
	err = os.WriteFile(fmt.Sprintf("%d-response.xml", n), body, 0777)
	if err != nil {
		panic(err)
	}
	fmt.Printf("response written to %d-response.xml\n", n)
	*/
	return nil
}
