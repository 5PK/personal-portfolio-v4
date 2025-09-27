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
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	ag "github.com/5PK/agility-fetch-go"

	"github.com/a-h/templ"
)

// Map CMS route IDs (or names) to templ components
var PageComponents = map[string]func() templ.Component{
	"help": partials.Help,
	"home": partials.Home,
}

// randomSelectImages randomly selects up to count unique images from the input slice
func randomSelectImages(images []types.ImageGalleryItem, count int) []types.ImageGalleryItem {

	if len(images) <= count {
		return images
	}

	// Create a random source
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create indices slice
	indices := make([]int, len(images))
	for i := range indices {
		indices[i] = i
	}

	// Shuffle indices
	r.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	// Take first 'count' indices and build result
	result := make([]types.ImageGalleryItem, count)
	for i := 0; i < count; i++ {
		result[i] = images[indices[i]]
	}

	return result
}

func DecodeFields[T any](item ag.HeadlessContentItem) (T, error) {
	var typed T
	b, err := json.Marshal(item.Fields)
	if err != nil {
		return typed, err
	}
	err = json.Unmarshal(b, &typed)
	return typed, err
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
		highlights, err := CastField[[]ag.HeadlessContentItem](fields, "textHighlightList")
		if err != nil {
			panic(err) // or handle gracefully
		}

		// parse the fields
		var highlightFields []types.HighlightFields
		for _, hl := range highlights {
			fields, _ := DecodeFields[types.HighlightFields](hl)
			highlightFields = append(highlightFields, fields)

		}

		// find max text length
		maxLen := 0
		for _, hl := range highlightFields {
			if len(hl.Text) > maxLen {
				maxLen = len(hl.Text)
			}
		}

		// build inline array with whitespace
		var highlightsWithWhiteSpace []types.CommandRow

		for _, hl := range highlightFields {
			padding := strings.Repeat(" ", maxLen-len(hl.Text)+4) // +4 for spacing
			highlightsWithWhiteSpace = append(highlightsWithWhiteSpace, types.CommandRow{
				Text: hl.Text,
				Ws:   padding,
				Desc: hl.Description,
			})
		}

		return components.CommandTitleWithDescription(
			fields["title"].(string),
			fields["description"].(string),
			highlightsWithWhiteSpace,
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
	"ImageGallery": func(fields map[string]interface{}) templ.Component {
		allImages, err := CastField[[]types.ImageGalleryItem](fields, "images")
		if err != nil {
			panic(err)
		}

		// Randomly select 6 unique images
		selectedImages := randomSelectImages(allImages, 6)

		header := fields["header"].(string)
		description := fields["description"].(string)

		return components.ImageGallery(header, description, selectedImages)
	},
}

var (
	sitemap     *map[string]ag.HeadlessContentSiteMapItem
	sitemapLock sync.RWMutex
)

var agilityClient *ag.APIClient

func InitializeAPI() {

	configuration := ag.NewConfiguration()
	configuration.AddDefaultHeader("APIKey", os.Getenv("AGILITY_API_KEY"))
	configuration.Servers = ag.ServerConfigurations{
		{
			URL:         "https://api.aglty.io",
			Description: "Agility CMS API Server",
		},
	}

	agilityClient = ag.NewAPIClient(configuration)

}

// Exported accessor so other packages can use it
func AgilityClient() *ag.APIClient {
	return agilityClient
}

// Refresh updates the sitemap from API
func RefreshSitemap() {
	sm := GetSitemapFlat() // <- your API call
	sitemapLock.Lock()
	sitemap = sm
	sitemapLock.Unlock()
}

// Get returns the current sitemap
func GetCurrentSitemap() *map[string]ag.HeadlessContentSiteMapItem {
	sitemapLock.RLock()
	defer sitemapLock.RUnlock()
	return sitemap
}

func RenderPage(ctx context.Context, w io.Writer, page ag.HeadlessContentPage) error {
	log.Println("here" + page.GetName())
	for _, modules := range page.Zones {
		for _, m := range modules {
			log.Println("module: " + m.GetModule())
			renderer, ok := PageComponentRenderer[m.GetModule()]
			if !ok {
				// fallback for unknown modules
				fmt.Fprintf(w, "<!-- unknown module: %s -->", m.GetModule())
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

// GetHomeCommandContent gets the CMS content for the home command
func GetHomeCommandContent() templ.Component {
	sm := GetCurrentSitemap()
	homeRoute := "/terminal/commands/home"

	for route, sitemapPage := range *sm {
		if route == homeRoute {
			page := GetPage(*sitemapPage.PageID)
			log.Println(json.Marshal(page))
			return createPageComponent(*page)
		}
	}

	// Fallback to static home if not found in CMS
	return partials.Home()
}

// createPageComponent converts a CMS page to a templ component
func createPageComponent(page ag.HeadlessContentPage) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return RenderPage(ctx, w, page)
	})
}

func GetPage(pageID int32) *ag.HeadlessContentPage {
	log.Println("page call" + strconv.Itoa(int(pageID)) + "234c2f44-u" + "en-us")
	resp, _, err  := agilityClient.PageAPI.PageIdGet(context.Background(), "234c2f44-u", ag.FETCH, "en-us", pageID).ContentLinkDepth(5).ExpandAllContentLinks(true).Execute()

	if err != nil {
		panic(err)
	}

	pageJson, _ := json.Marshal(resp)
	log.Println(string(pageJson))
	return resp
}

func GetSitemapFlat() *map[string]ag.HeadlessContentSiteMapItem {

	// agilityClient.SitemapAPI.SitemapFlatChannelNameGet(context.Background(), )
	resp, _, err := agilityClient.SitemapAPI.SitemapFlatChannelNameGet(context.Background(), "234c2f44-u", ag.FETCH, "en-us", "website").Execute()

	if err != nil {
		panic(err)
	}

	return resp
}
