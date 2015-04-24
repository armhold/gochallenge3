package gochallenge3
import "encoding/json"

type InstagramPagination  struct {
	MaxTagID    string `json:"max_tag_id"`
	MinTagID    string `json:"min_tag_id"`
	NextURL     string `json:"next_url"`
}

type InstagramImage struct {
	Url    string  `json:"url"`
	Width  int32   `json:"width"`
	Height int32   `json:"height"`
}

type InstagramImageSet struct {
	LowRes      InstagramImage `json:"low_resolution"`
	Thumb       InstagramImage `json:"thumbnail"`
	StandardRes InstagramImage `json:"standard_resolution"`
}

type InstagramData struct {
	Tags []string                  `json:"tags"`
	ImageSet InstagramImageSet     `json:"images"`
}

type InstagramResponse struct {
	Pagination InstagramPagination `json:"pagination"`
	Data       []InstagramData     `json:"data"`
}


func ParseInstagramJSON(jsonBytes []byte) (InstagramResponse, error) {
	var result InstagramResponse

	err := json.Unmarshal(jsonBytes, &result)
	return result, err
}
