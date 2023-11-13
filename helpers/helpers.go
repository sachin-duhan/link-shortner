package helpers

import (
	"log"
	"os"
	"strings"
)

func EnforceHTTP(url string) string {
	// make sure that every URL should start with HTTP protocol.
	if url[:4] != "http" {
		return "http:/" + url
	}
	return url
}

func RemoveDomainError(url string) bool {
	// basically this functions removes all the commonly found
	// prefixes from URL such as http, https, www
	// then checks of the remaining string is the DOMAIN itself

	DOMAIN := os.Getenv("DOMAIN")

	if len(DOMAIN) == 0 {
		log.Fatal("Missing domain environment variable.")
		panic("Missing domain name in env.")
	}

	if url == DOMAIN {
		return false
	}

	url = strings.Replace(url, "http://", "", 1)
	url = strings.Replace(url, "https://", "", 1)
	url = strings.Replace(url, "www.", "", 1)
	url = strings.Split(url, "/")[0]

	return url != DOMAIN
}