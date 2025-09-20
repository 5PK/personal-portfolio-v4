package types

type HighlightFields struct {
	Text        string `json:"text"`
	Description string `json:"description"`
}

type TextHighlightWithDescription = ContentItem[HighlightFields]
