package agility

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kevin-portfolio/internal/agility_cms/types"
	"kevin-portfolio/views/components"
	"kevin-portfolio/views/partials"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/a-h/templ"
)

// Map CMS route IDs (or names) to templ components
var PageComponents = map[string]func() templ.Component{
	"help": partials.Help,
	"home": partials.Home,
}

func CastField[T any](fields map[string]any, key string) (T, error) {
	var out T
	raw, ok := fields[key]
	if !ok {
		return out, fmt.Errorf("field %q not found", key)
	}

	// Re-marshal then unmarshal to the target type
	b, err := json.Marshal(raw)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return out, err
	}
	return out, nil
}

var PageComponentRenderer = map[string]func(fields map[string]interface{}) templ.Component{
	"Links": func(fields map[string]interface{}) templ.Component {

		githubRaw := fields["github"]
		githubBytes, _ := json.Marshal(githubRaw) // re-marshal to JSON
		var github types.LinkField
		_ = json.Unmarshal(githubBytes, &github)

		linkedInRaw := fields["linkedin"]
		linkedInBytes, _ := json.Marshal(linkedInRaw) // re-marshal to JSON
		var linkedIn types.LinkField
		_ = json.Unmarshal(linkedInBytes, &linkedIn)

		return components.Links(
			github,
			linkedIn,
		)
	},
	"CommandTitleWithDescription": func(fields map[string]interface{}) templ.Component {
		// retrieve the list of text highlights
		highlights, err := CastField[[]types.TextHighlightWithDescription](fields, "textHighlightList")
		if err != nil {
			panic(err) // or handle gracefully
		}

		return components.CommandTitleWithDescription(
			fields["title"].(string),
			fields["description"].(string),
			highlights,
			fields["footnote"].(string),
		)
	},
	"NameHeaderQuote": func(fields map[string]interface{}) templ.Component {
		return components.NameHeaderQuote(
			fields["headerName"].(string),
			fields["jobTitle"].(string),
			fields["quote"].(string),
			fields["quoteAuthor"].(string),
			fields["quoteSource"].(string),
		)
	},
	"TextBlockWithHeader": func(fields map[string]interface{}) templ.Component {
		textBlockString := fields["textBlock"].(string)
		return components.TextBlockWithHeader(
			fields["header"].(string),
			strings.ReplaceAll(textBlockString, "\r\n", "<br>"),
		)
	},
}

var (
	sitemap     types.Sitemap
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
func GetCurrentSitemap() types.Sitemap {
	sitemapLock.RLock()
	defer sitemapLock.RUnlock()
	return sitemap
}

func RenderPage(ctx context.Context, w io.Writer, page types.Page) error {
	log.Println("here" + page.Name)
	for _, modules := range page.Zones {
		for _, m := range modules {
			log.Println("module: " + m.Module)
			renderer, ok := PageComponentRenderer[m.Module]
			if !ok {
				// fallback for unknown modules
				fmt.Fprintf(w, "<!-- unknown module: %s -->", m.Module)
				continue
			}

			component := renderer(m.Item.Fields)
			if err := component.Render(ctx, w); err != nil {
				return err
			}
		}
	}
	return nil
}


func GetPage(pageID int) types.Page {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.aglty.io/234c2f44-u/fetch/en-us/page/"+strconv.Itoa(pageID)+"?contentLinkDepth=2&expandAllContentLinks=true", nil)
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

	var page types.Page

	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		panic(err)
	}

	return page
}

// GetHomeCommandContent gets the CMS content for the home command
func GetHomeCommandContent() templ.Component {
	sm := GetCurrentSitemap()
	homeRoute := "/terminal/commands/home"

	for route, sitemapPage := range sm {
		if route == homeRoute {
			page := GetPage(sitemapPage.PageID)
			return createPageComponent(page)
		}
	}

	// Fallback to static home if not found in CMS
	return partials.Home()
}

// createPageComponent converts a CMS page to a templ component
func createPageComponent(page types.Page) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return RenderPage(ctx, w, page)
	})
}

func GetSitemapFlat() types.Sitemap {
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

	var sitemap types.Sitemap

	if err := json.NewDecoder(resp.Body).Decode(&sitemap); err != nil {
		panic(err)
	}

	return sitemap
}
