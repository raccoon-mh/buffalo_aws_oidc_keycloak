package actions

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gobuffalo/buffalo"
)

func MainHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("main/index.plush.html"))
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

func LoginHandler(c buffalo.Context) error {
	if c.Request().Method == "POST" {
		id := c.Param("id")
		password := c.Param("password")
		keycloakHost := os.Getenv("keycloakHost")
		realm := os.Getenv("realm")
		client := os.Getenv("client")
		clientSecret := os.Getenv("clientSecret")
		// RoleArn := os.Getenv("RoleArn")

		keycloakUrl := "https://" + keycloakHost + "/realms/" + realm + "/protocol/openid-connect/token"

		bodyData := url.Values{
			"username":      {id},
			"password":      {password},
			"client_id":     {client},
			"client_secret": {clientSecret},
			"grant_type":    {"password"},
		}

		req, _ := http.NewRequest("POST", keycloakUrl, bytes.NewBufferString(bodyData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		httpClient := &http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var tokenresponse TokenResponse
		jsonerr := json.Unmarshal(body, &tokenresponse)
		if jsonerr != nil {
			fmt.Println("Error parsing JSON:", err)
		}
		c.Session().Set("access_token", tokenresponse.AccessToken)

		return c.Redirect(302, "/user/home")
	}
	return c.Render(http.StatusOK, r.HTML("auth/index.plush.html"))
}

func LogoutHandler(c buffalo.Context) error {
	c.Session().Clear()
	return c.Redirect(302, "/")
}

type AssumeRoleWithWebIdentityResponse struct {
	AssumeRoleWithWebIdentityResult AssumeRoleWithWebIdentityResult `xml:"AssumeRoleWithWebIdentityResult"`
	ResponseMetadata                ResponseMetadata                `xml:"ResponseMetadata"`
}

type AssumeRoleWithWebIdentityResult struct {
	Audience               string          `xml:"Audience"`
	AssumedRoleUser        AssumedRoleUser `xml:"AssumedRoleUser"`
	Provider               string          `xml:"Provider"`
	Credentials            Credentials     `xml:"Credentials"`
	SubjectFromWebIdentity string          `xml:"SubjectFromWebIdentityToken"`
}

type AssumedRoleUser struct {
	AssumedRoleId string `xml:"AssumedRoleId"`
	Arn           string `xml:"Arn"`
}

type Credentials struct {
	AccessKeyId     string `xml:"AccessKeyId"`
	SecretAccessKey string `xml:"SecretAccessKey"`
	SessionToken    string `xml:"SessionToken"`
	Expiration      string `xml:"Expiration"`
}

type ResponseMetadata struct {
	RequestId string `xml:"RequestId"`
}

func GetStsTokenHandler(c buffalo.Context) error {

	RoleArn := os.Getenv("RoleArn")
	stsurl := "https://sts.amazonaws.com"

	paramData := url.Values{
		"DurationSeconds":  {"900"},
		"Action":           {"AssumeRoleWithWebIdentity"},
		"Version":          {"2011-06-15"},
		"RoleSessionName":  {"web-identity-federation"},
		"RoleArn":          {RoleArn},
		"WebIdentityToken": {c.Session().Get("access_token").(string)},
	}

	stsurl += "?" + paramData.Encode()

	req, _ := http.NewRequest("GET", stsurl, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var stsxml AssumeRoleWithWebIdentityResponse
	xmlerr := xml.Unmarshal(body, &stsxml)
	if xmlerr != nil {
		fmt.Println("Error parsing XML:", err)
	}

	c.Set("stsxml", stsxml)

	return c.Render(http.StatusOK, r.HTML("main/index.plush.html"))
}
