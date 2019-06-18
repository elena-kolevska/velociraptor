package twilio

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"sort"
)

func IsValidTwilioSignature(scheme string, host string, authToken string, URL string, postForm url.Values, xTwilioSignature string) bool {
	return xTwilioSignature == getExpectedTwilioSignature(scheme, host, authToken, URL, postForm)
}

func getExpectedTwilioSignature(scheme string, host string, authToken string, URL string, postForm url.Values) string {
	// Take the full URL of the request URL
	str := scheme + "://" + host + URL

	// If the request is a POST, sort all of the POST parameters alphabetically
	keys := make([]string, 0, len(postForm))
	for key := range postForm {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Iterate through the sorted list of POST parameters and append
	// the variable name and value to the end of the URL string.
	for _, key := range keys {
		str += key + postForm[key][0]
	}

	// Sign the resulting string with HMAC-SHA1 using Twilio AuthToken as the key
	mac := hmac.New(sha1.New, []byte(authToken))
	mac.Write([]byte(str))
	expectedMac := mac.Sum(nil)

	// Base64 encode the resulting hash value.
	return base64.StdEncoding.EncodeToString(expectedMac)
}
