package gochallenge3

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// simple wrapper for Instagram REST API
type ImageSource interface {
	Search(s string) ([]ImageURL, error)
}

type InstagramImageSource struct {
	clientID string
}

func NewInstagramImageSource(clientID string) *InstagramImageSource {
	return &InstagramImageSource{clientID: clientID}
}

func (i *InstagramImageSource) Search(s string) ([]ImageURL, error) {
	log.Println("starting search...")
	u, err := i.instagramAPIUrl(s)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("image fetch failed: %v", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	got, err := ParseInstagramJSON([]byte(body))

	if err != nil {
		return nil, err
	}

	var result []ImageURL

	for _, datum := range got.Data {
		url := datum.ImageSet.Thumb.Url
		imageURL := ImageURL(url)
		result = append(result, imageURL)
	}

	log.Println("search complete.	")

	return result, nil
}

func (i *InstagramImageSource) instagramAPIUrl(searchTag string) (string, error) {
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
	parameters.Add("client_id", i.clientID)
	parameters.Add("count", "50")
	u.RawQuery = parameters.Encode()

	return u.String(), nil
}
