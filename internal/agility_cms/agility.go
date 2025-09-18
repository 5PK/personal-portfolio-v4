package agility

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kevin-portfolio/views/partials"
	"net/http"
	"os"
	"sync"

	"github.com/a-h/templ"
)

type Page struct {
	ID    string                      `json:"id"`
	Zones map[string][]ModuleInstance `json:"zones"`
}

type ModuleInstance struct {
	Module string                 `json:"module"`
	Fields map[string]interface{} `json:"fields"`
}

type Sitemap map[string]Page

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

// Map CMS route IDs (or names) to templ components
var PageComponents = map[string]func() templ.Component{
	"help": partials.Help,
	"home": partials.Home,
}

var PageComponentRenderer = map[string]func(fields map[string]interface{}) templ.Component{
	"TextBlockWithHeader": func(fields map[string]interface{}) templ.Component {
		return partials.Home(
		// fields["header"].(string),
		// fields["text"].(string),
		)
	},
}

var (
	sitemap     Sitemap
	sitemapLock sync.RWMutex
)

// Refresh updates the sitemap from API
func RefreshSitemap() {
	sm := GetSitemapFlat() // <- your API call
	sitemapLock.Lock()
	sitemap = sm
	sitemapLock.Unlock()
}

// Get returns the current sitemap
func GetCurrentSitemap() Sitemap {
	sitemapLock.RLock()
	defer sitemapLock.RUnlock()
	return sitemap
}

func GetSitemapFlat() Sitemap {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.aglty.io/234c2f44-u/fetch/en-us/sitemap/flat/website", nil)
	if err != nil {
		panic(err)
	}

	// Add headers
	// req.Header.Set("Authorization", "Bearer "+os.Getenv("API_KEY"))
	req.Header.Set("apikey", "defaultlive.57b62aceeb699b8085def6525f24a0ce81584ecd31929c523e716d4187beea06")
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("API error: %s\n%s", resp.Status, string(body)))
	}

	var sitemap Sitemap

	if err := json.NewDecoder(resp.Body).Decode(&sitemap); err != nil {
		panic(err)
	}

	return sitemap
}

func RenderPage(ctx context.Context, w io.Writer, page Page) error {
	for _, modules := range page.Zones {
		for _, m := range modules {
			renderer, ok := PageComponentRenderer[m.Module]
			if !ok {
				// fallback for unknown modules
				fmt.Fprintf(w, "<!-- unknown module: %s -->", m.Module)
				continue
			}

			component := renderer(m.Fields)
			if err := component.Render(ctx, w); err != nil {
				return err
			}
		}
	}
	return nil
}
