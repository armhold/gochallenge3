package gochallenge3

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	searchTag, err := UrlEncode(searchTag)
	if err != nil {
		return "", err
	}

	u, err := url.Parse("https://api.instagram.com/v1/tags")
	if err != nil {
		return "", err
	}

	u.Path += "/" + searchTag
	u.Path += "/media/recent"
	parameters := url.Values{}
	parameters.Add("client_id", i.ClientID)
	parameters.Add("count", "50")
	u.RawQuery = parameters.Encode()

	return u.String(), nil
}
