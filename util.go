package gochallenge3
import (
    "net/url"
    "regexp"
)

var (
    prohibited *regexp.Regexp
)


func init() {
    var err error

    // disallow ".." and "/" strings to appear in the search term
    prohibited, err = regexp.Compile("\\.\\.|/")
    if err != nil {
        panic(err)
    }
}

// UrlEncode encodes a string like Javascript's encodeURIComponent(), but also strips slashes and ".."
func UrlEncode(s string) (string, error) {
    s = prohibited.ReplaceAllString(s, "")

    u, err := url.Parse(s)
    if err != nil {
        return "", err
    }
    return u.String(), nil
}
