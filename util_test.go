package gochallenge3

import (
    "testing"
)

func TestUrlEncodeSpaces(t *testing.T) {
    input := "foo bar"
    got, err := UrlEncode(input)
    want := "foo%20bar"

    if err != nil {
        t.Fatal(err)
    }

    if got != want {
        t.Errorf("UrlEncoded(%q) => %q, want %q", input, got, want)
    }
}

func TestUrlEncodeDotsAndSlashes(t *testing.T) {
    input := "foo/../etc/password"
    got, err := UrlEncode(input)
    want := "fooetcpassword"

    if err != nil {
        t.Fatal(err)
    }

    if got != want {
        t.Errorf("UrlEncoded(%q) => %q, want %q", input, got, want)
    }
}
