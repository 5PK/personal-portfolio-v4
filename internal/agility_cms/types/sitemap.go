package types

type Sitemap map[string]SitemapPage

type SitemapPage struct {
	Title    *string `json:"title"` // nullable
	Name     string  `json:"name"`
	PageID   int     `json:"pageID"`
	MenuText string  `json:"menuText"`
	Visible  struct {
		Menu    bool `json:"menu"`
		Sitemap bool `json:"sitemap"`
	} `json:"visible"`
	Path     string      `json:"path"`
	Redirect interface{} `json:"redirect"` // could be string or null
	IsFolder bool        `json:"isFolder"`
}
