package types

type Page struct {
	PageID           int                   `json:"pageID"`
	Name             string                `json:"name"`
	Path             *string               `json:"path"`
	Title            string                `json:"title"`
	MenuText         string                `json:"menuText"`
	PageType         string                `json:"pageType"`
	TemplateName     string                `json:"templateName"`
	RedirectURL      string                `json:"redirectUrl"`
	SecurePage       bool                  `json:"securePage"`
	ExcludeFromCache bool                  `json:"excludeFromOutputCache"`
	Visible          Visible               `json:"visible"`
	SEO              SEO                   `json:"seo"`
	Scripts          Scripts               `json:"scripts"`
	Properties       Properties            `json:"properties"`
	Zones            map[string][]ZoneItem `json:"zones"`
}

type Visible struct {
	Menu    bool `json:"menu"`
	Sitemap bool `json:"sitemap"`
}

type SEO struct {
	MetaDescription string `json:"metaDescription"`
	MetaKeywords    string `json:"metaKeywords"`
	MetaHTML        string `json:"metaHTML"`
	MenuVisible     *bool  `json:"menuVisible"`
	SitemapVisible  *bool  `json:"sitemapVisible"`
}

type Scripts struct {
	ExcludedFromGlobal bool    `json:"excludedFromGlobal"`
	Top                *string `json:"top"`
	Bottom             *string `json:"bottom"`
}

type Properties struct {
	State     int    `json:"state"`
	Modified  string `json:"modified"` // keep as string for now
	VersionID int    `json:"versionID"`
}

type ZoneItem struct {
	Module string `json:"module"`
	Item   Item   `json:"item"`
}

type Item struct {
	ContentID  int            `json:"contentID"`
	Properties ItemProperties `json:"properties"`
	Fields     map[string]any `json:"fields"`
}

type ItemProperties struct {
	State          int    `json:"state"`
	Modified       string `json:"modified"` // keep as string for now
	VersionID      int    `json:"versionID"`
	ReferenceName  string `json:"referenceName"`
	DefinitionName string `json:"definitionName"`
	ItemOrder      int    `json:"itemOrder"`
}
