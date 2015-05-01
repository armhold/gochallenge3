package gochallenge3

type ImageURL string

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
