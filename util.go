package gochallenge3
import (
    "net/url"
    "regexp"
)

var (
    dotRegexp *regexp.Regexp
    slashRegexp *regexp.Regexp
)


func init() {
    var err error

    dotRegexp, err = regexp.Compile("\\.\\.")
    if err != nil {
        panic(err)
    }

    slashRegexp, err = regexp.Compile("/")
    if err != nil {
        panic(err)
    }
}

// UrlEncode encodes a string like Javascript's encodeURIComponent(), but also strips slashes and ".."
func UrlEncode(s string) (string, error) {
    s = dotRegexp.ReplaceAllString(s, "")
    s = slashRegexp.ReplaceAllString(s, "")

    u, err := url.Parse(s)
    if err != nil {
        return "", err
    }
    return u.String(), nil
}
