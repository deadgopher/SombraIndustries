package controller

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type iEve interface {
	esiLink() string
	MakeRequest() ([]byte, error)
	RegisterEve(*gin.RouterGroup)
}

type eve struct {
	eveRoot
	base            string
	callback        string
	scope           iScope
	state           string
	client          string
	postDestination string
	creds           string
}

func (x *eve) RegisterEve(r *gin.RouterGroup) {
	r.GET("/", x.handleCallback)
}

func newEveController(r eveRoot) *eve {
	x := &eve{
		eveRoot: r,
	}
	x.base = "https://login.eveonline.com/v2/oauth/authorize/?response_type=code&"
	x.callback = "redirect_uri=" + url.QueryEscape("https://localhost:8080/api/eve/")
	x.scope = scope{}.new()
	x.state = "state=superSecretEveApiState"
	x.client = fmt.Sprintf("client_id=%v", os.Getenv("EVE_CLIENT_ID"))
	x.postDestination = "https://login.eveonline.com/v2/oauth/token"
	return x
}
func (x *eve) esiLink() string {
	return fmt.Sprintf(x.base+"%v&%v&%v&%v", x.callback, x.client, x.scope.String(), x.state)
}

func (x *eve) verify(tok string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", "https://login.eveonline.com/oauth/verify", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (x *eve) handleCallback(c *gin.Context) {
	x.creds = "grant_type=authorization_code&code=" + c.Query("code")
	data := fmt.Sprintf("%v:%v", os.Getenv("EVE_CLIENT_ID"), os.Getenv("EVE_SECRET_KEY"))
	encData := b64.URLEncoding.EncodeToString([]byte(data))
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", x.postDestination, strings.NewReader(x.creds))
	if err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	req.Header.Set("Authorization", "Basic "+encData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "login.eveonline.com")

	res, err := client.Do(req)
	if err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		respond(c, 500, false, foo(err))
		return
	}

	pilot, err := x.PilotFromToken(body)
	if err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	r, err := x.verify(pilot.AccessToken)
	if err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	var bar struct {
		CharacterID string
	}
	if err := json.Unmarshal(r, &bar); err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	pilot.PilotID = bar.CharacterID
	if err := pilot.Save(); err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	x.Show(c.Writer, "index")
}
