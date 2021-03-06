package main

import (
	"net/http"
	"fmt"
	"crypto/rand"
	"encoding/base64"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Example: this will give us a 44 byte, base64 encoded output
	token, err := GenerateRandomString(32)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprintf(w, token)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":80", nil)
}

// Thanks to https://stackoverflow.com/a/32351471/795407

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	outputBytes := make([]byte, n)
	_, err := rand.Read(outputBytes)
	// Note that err == nil only if we read len(outputBytes) bytes.
	if err != nil {
		return nil, err
	}

	return outputBytes, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(s int) (string, error) {
	randomBytes, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(randomBytes), err
}