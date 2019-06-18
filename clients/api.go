package clients

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ModelLifeClient struct
type ModelLifeClient struct {
	BaseURL               string
	GetTokenURL           string
	UpdateConversationURL string
	ClientID              string
	ClientSecret          string
	AccessToken           string
	RefreshToken          string
	HttpClient            http.Client
}

// GetAccessToken sends a PUT request to the commune X API,
// updating the last activity timestamp of a certain conversation
func (c *ModelLifeClient) GetAccessToken() (err error) {
	urlStr := c.BaseURL + c.GetTokenURL

	data := url.Values{}
	data.Add("grant_type", "client_credentials")
	data.Add("client_id", c.ClientID)
	data.Add("client_secret", c.ClientSecret)

	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := c.HttpClient.Do(r)
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject TokenResponse
	if err := json.Unmarshal(responseData, &responseObject); err != nil {
		log.Fatal(err)
		return err
	}

	token := responseObject.Data
	c.AccessToken = token.AccessToken
	c.RefreshToken = token.RefreshToken

	return nil
}

//UpdateLastActivity sends a request to the API to update a conversations's updated_at timestamp
func (c ModelLifeClient) UpdateLastActivity(channelID *string, timestamp *int64) (err error) {
	urlStr := c.BaseURL + c.UpdateConversationURL + "/" + *channelID

	data := url.Values{}
	data.Add("timestamp", strconv.FormatInt(*timestamp, 10))

	r, _ := http.NewRequest("PUT", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Authorization", "Bearer "+c.AccessToken)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := c.HttpClient.Do(r)

	if resp.StatusCode != 200 {
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var responseObject GenericResponse
		if err := json.Unmarshal(responseData, &responseObject); err != nil {
			return err
		}
		return errors.New(responseObject.Errors.Title)
	}

	return nil
}

// A TokenResponse struct to map the Entire API Response
type TokenResponse struct {
	Data   Token    `json:"data,omitempty"`
	Errors APIError `json:"errors,omitempty"`
}

// A GenericResponse struct to map the Entire API Response
type GenericResponse struct {
	Errors APIError `json:"errors,omitempty"`
}

// An APIError struct to map errors from the API
type APIError struct {
	Title string `json:"title"`
}

// A Token struct to map the access token from the API
type Token struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int32  `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
