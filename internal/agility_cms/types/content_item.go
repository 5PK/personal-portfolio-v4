package types

type ContentItem[T any] struct {
	ContentID  int                    `json:"contentID"`
	Properties ContentItemProperties  `json:"properties"`
	Fields     T 					  `json:"fields"` // dynamic fields
	SEO        interface{}            `json:"seo"`    // can refine if you know schema
}

type ContentItemProperties struct {
	State          int    `json:"state"`
	Modified       string `json:"modified"`
	VersionID      int    `json:"versionID"`
	ReferenceName  string `json:"referenceName"`
	DefinitionName string `json:"definitionName"`
	ItemOrder      int    `json:"itemOrder"`
}
