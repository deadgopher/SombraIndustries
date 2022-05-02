package controller

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type eveController struct {
	base             string
	callback         string
	scope            string
	state            string
	client           string
	post_destination string
	response         *eveResponse
}

func newEveController(r root) *eveController {
	x := &eveController{
		"https://login.eveonline.com/v2/oauth/authorize/?response_type=code&",
		"redirect_uri=" + url.QueryEscape("https://localhost:8080/api/eve/"),
		"scope=esi-characters.read_blueprints.v1 esi-corporations.read_contacts.v1",
		"state=superSecretEveApiState",
		"client_id=" + os.Getenv("EVE_CLIENT_ID"),
		"https://login.eveonline.com/v2/oauth/token",
		nil,
	}
	return x
}
func (x *eveController) Link() string {
	return fmt.Sprintf(x.base+"%v&%v&%v&%v", x.callback, x.client, x.scope, x.state)
}

func (x *eveController) handleCallback(c *gin.Context) {
	code := c.Query("code")
	y := url.Values{}
	y.Set("grant_type", "authorization_code&code="+code)

	data := fmt.Sprintf("%v:%v", os.Getenv("EVE_CLIENT_ID"), os.Getenv("EVE_SECRET_KEY"))
	encData := b64.URLEncoding.EncodeToString([]byte(data))

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", x.post_destination, strings.NewReader(y.Encode()))
	if err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	req.Header.Set("Authorization", "Basic "+encData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "login.eveonline.com")

	res, err2 := client.Do(req)
	if err2 != nil {
		respond(c, 500, false, foo(err))
		return
	}
	defer res.Body.Close()

	// expected json payload
	// {
	// 	"access_token": <JWT token>,
	// 	"expires_in": 1199,
	// 	"token_type": "Bearer",
	// 	"refresh_token": <unique string>
	// }
	var decodedRes *eveResponse
	if err := json.NewDecoder(res.Body).Decode(&decodedRes); err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	x.response = decodedRes
	respond(c, 200, true, decodedRes)
}

type eveResponse struct {
	AccessToken  string    `json:"access_token"`
	ExpiresIn    time.Time `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
}

func (x *eveController) registerRoutes(r *gin.RouterGroup) {
	r.GET("/", x.handleCallback)
}

// func call(url, method string) error {
//     client := &http.Client{
//         Timeout: time.Second * 10,
//     }
//     req, err := http.NewRequest(method, url, nil)
//     if err != nil {
//         return fmt.Errorf("Got error %s", err.Error())
//     }
//     req.Header.Set("user-agent", "golang application")
//     req.Header.Add("foo", "bar1")
//     req.Header.Add("foo", "bar2")
//     response, err := client.Do(req)
//     if err != nil {
//         return fmt.Errorf("Got error %s", err.Error())
//     }
//     defer response.Body.Close()
//     return nil
// }

// Create a Base64 encoded string, including padding, where the contents before encoding are your application’s client ID, followed by a :, followed by your application’s secret key (e.g. Base64(<client_id>:<secret_key>)). For example, given the input CLIENT_ID:CLIENT_SECRET, the resulting string should be Q0xJRU5UX0lEOkNMSUVOVF9TRUNSRVQ=.

// You will need to send the following HTTP headers (replace anything between <>, including <>):
// Authorization: Basic <Base64 encoded credentials>
// Content-Type: application/x-www-form-urlencoded
// Host: login.eveonline.com
