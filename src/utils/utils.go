package utils

import (
	"encoding/base64"

	"golang.org/x/text/encoding/unicode"
)

// Backspace character
var Backspace = []byte{8}

// Newline character
var Newline = []byte{10}

// Nullbyte character
var Nullbyte = []byte{0}

// SliceStringContains check if a string is present in a slice of string
func SliceStringContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// SliceByteContains check if a byte is present in a slice of byte
func SliceByteContains(s []byte, e byte) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Min returns the minimal integer between two integer
func Min(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// Iface represents an interface on the system with its name and its IP
type Iface struct {
	Name string
	IP   string
}

func Utf16leBase64(s string) (string, error) {
	var stringB64 = ""
	utfEncoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	ut16LeEncodedMessage, err := utfEncoder.String(s)
	if err == nil {
		stringB64 = base64.StdEncoding.EncodeToString([]byte(ut16LeEncodedMessage))
	}
	return stringB64, err
}
