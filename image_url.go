package gochallenge3
import (
    "os"
    "io"
    "bufio"
)

type ImageURL string


func ToFile(urls []ImageURL, filename string) error {
    file, err := os.Create(filename)

    if err != nil {
        return err
    }
    defer file.Close()

    for _, url := range urls {
        line := string(url) + "\n"
        _, err := io.WriteString(file, line)
        if err != nil {
            return err
        }
    }

    return nil
}

func FromFile(filename string) ([]ImageURL, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var result []ImageURL
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        url := ImageURL(line)
        result = append(result, url)
    }

    return result, scanner.Err()
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
