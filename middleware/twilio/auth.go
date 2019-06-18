package twilio

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"github.com/labstack/echo"
	"net/http"
	"net/url"
	"sort"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if releaseStage != "localhost" {
			scheme := "http"
			if c.IsTLS() {
				scheme = "https"
			}

			formParameters, err := c.FormParams()
			if err != nil {
				return c.String(http.StatusUnauthorized, "Something's wrong with your request")
			}

			if twilio.isValidTwilioSignature(scheme, c.Request().Host, twilioAPIToken, c.Request().RequestURI, formParameters, c.Request().Header.Get("X-Twilio-Signature")) == false {
				return c.String(http.StatusUnauthorized, "You're not Twilio")
			}
		}

		return next(c)
	}
}

func isValidTwilioSignature(scheme string, host string, authToken string, URL string, postForm url.Values, xTwilioSignature string) bool {
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
