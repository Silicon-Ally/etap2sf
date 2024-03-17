package client

import (
	"github.com/tzmfreedom/go-metaforce"
)

func (c *Client) ListLayouts() (*metaforce.ListMetadataResponse, error) {
	return c.gc.MetadataClient.ListMetadata([]*metaforce.ListMetadataQuery{{
		Type: "Layout",
	}})
}

func (c *Client) ListCompactLayout() (*metaforce.ListMetadataResponse, error) {
	return c.gc.MetadataClient.ListMetadata([]*metaforce.ListMetadataQuery{{
		Type: "CompactLayout",
	}})
}

func (c *Client) ListFlexiPages() (*metaforce.ListMetadataResponse, error) {
	return c.gc.MetadataClient.ListMetadata([]*metaforce.ListMetadataQuery{{
		Type: "FlexiPage",
	}})
}
