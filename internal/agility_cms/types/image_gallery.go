package types

// ImageGalleryItem represents an individual image in the gallery
type ImageGalleryItem struct {
	Label       *string `json:"label"`
	URL         string  `json:"url"`
	Target      *string `json:"target"`
	Filesize    int     `json:"filesize"`
	PixelHeight string  `json:"pixelHeight"`
	PixelWidth  string  `json:"pixelWidth"`
	Height      int     `json:"height"`
	Width       int     `json:"width"`
}

// ImageGalleryFields represents the fields for an ImageGallery module
type ImageGalleryFields struct {
	Header      string             `json:"header"`
	Description string             `json:"description"`
	Images      []ImageGalleryItem `json:"images"`
}