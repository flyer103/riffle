package riffle

import (
	"encoding/xml"
	"fmt"
	"os"
)

// OPML represents the root OPML structure
type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

// Head represents the OPML head section
type Head struct {
	Title string `xml:"title"`
}

// Body represents the OPML body section
type Body struct {
	Outlines []Outline `xml:"outline"`
}

// Outline represents an OPML outline element
type Outline struct {
	Title    string    `xml:"title,attr"`
	Text     string    `xml:"text,attr"`
	Type     string    `xml:"type,attr"`
	XMLURL   string    `xml:"xmlUrl,attr"`
	HTMLURL  string    `xml:"htmlUrl,attr"`
	Outlines []Outline `xml:"outline"`
}

// Feed represents a feed from OPML
type Feed struct {
	Title string
	URL   string
}

// ParseOPML parses an OPML file and returns a list of feeds
func ParseOPML(filename string) ([]Feed, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read OPML file: %w", err)
	}

	var doc OPML
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse OPML file: %w", err)
	}

	var feeds []Feed
	for _, outline := range doc.Body.Outlines {
		feeds = append(feeds, extractFeeds(outline)...)
	}

	return feeds, nil
}

// extractFeeds recursively extracts feeds from an outline and its children
func extractFeeds(outline Outline) []Feed {
	var feeds []Feed

	// If this outline is a feed
	if outline.XMLURL != "" {
		title := outline.Title
		if title == "" {
			title = outline.Text
		}
		feeds = append(feeds, Feed{
			Title: title,
			URL:   outline.XMLURL,
		})
	}

	// Process child outlines
	for _, child := range outline.Outlines {
		feeds = append(feeds, extractFeeds(child)...)
	}

	return feeds
}
