package gochallenge3

import (
	"regexp"
	"strings"
	"fmt"
)

type ImageURL string

var (
	imageExtensionRegexp *regexp.Regexp
)

func init() {
	imageExtensionRegexp = regexp.MustCompile(`http.*(\.gif|\.jpg|\.png)\?.*`)
}


func ToRows(rowLen int, imageURLs []ImageURL) [][]ImageURL {
	i := 0

	var result [][]ImageURL

	for i < len(imageURLs) {
		end := i + rowLen
		if end > len(imageURLs) {
			end = len(imageURLs)
		}

		result = append(result, imageURLs[i:end])
		i = end
	}

	return result
}

func (i ImageURL) guessImageExtension() (string, error) {
	lower := strings.ToLower(string(i))
	match := imageExtensionRegexp.FindStringSubmatch(lower)
	if len(match) == 2 {
		return match[1], nil
	}

	return "", fmt.Errorf("failed to guess image extension: %s", i)
}
