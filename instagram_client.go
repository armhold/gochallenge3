package gosaic

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

var (
	// disallow "..", "/" and " " strings to appear in the search term
	prohibited = regexp.MustCompile(`\.\.|/|\s+`)
)

// simple client for the Instagram Search API
type InstagramClient struct {
	ClientID string
}

func (i *InstagramClient) Search(s string, maxResults int) ([]ImageURL, error) {
	log.Println("starting search...")
	instagramURL, err := i.instagramAPIUrl(s)
	if err != nil {
		return nil, err
	}

	var result []ImageURL

	for instagramURL != "" && len(result) < maxResults {
		log.Printf("searching URL: %s", instagramURL)

		imageURLs, nextURL, err := i.searchPaginated(instagramURL)
		if err != nil {
			return nil, err
		}

		for _, url := range imageURLs {
			result = append(result, url)
		}

		instagramURL = nextURL
	}

	log.Println("search complete.")

	return result[0:maxResults], nil
}

func (i *InstagramClient) searchPaginated(instagramURL string) (imageURLs []ImageURL, nextUrl string, err error) {
	req, err := http.NewRequest("GET", instagramURL, nil)
	if err != nil {
		return nil, "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.New(fmt.Sprintf("image search failed, http error: %v", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	got, err := ParseInstagramJSON([]byte(body))

	if err != nil {
		return nil, "", err
	}

	for _, datum := range got.Data {
		url := datum.ImageSet.Thumb.Url
		imageURL := ImageURL(url)
		imageURLs = append(imageURLs, imageURL)
	}

	return imageURLs, got.Pagination.NextURL, nil
}

func (i *InstagramClient) instagramAPIUrl(searchTag string) (string, error) {

	// NB: can't convert spaces to "+" or "%20"; the API just flat out rejects them,
	// as well as %2F for slashes, etc.
	cleanedTag := prohibited.ReplaceAllString(searchTag, "")

	u, err := url.Parse(fmt.Sprintf("https://api.instagram.com/v1/tags/%s/media/recent", cleanedTag))
	if err != nil {
		return "", err
	}

	parameters := url.Values{}
	parameters.Add("client_id", i.ClientID)
	parameters.Add("count", "50")
	u.RawQuery = parameters.Encode()

	return u.String(), nil
}
