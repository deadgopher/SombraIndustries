package controller

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"germ/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// sharing
type iESIController interface {
	register(*gin.RouterGroup)
	getLink() string
}

// requires
type iESIRoot interface {
	CreateToken(data string, exp time.Duration) (string, error)
	SetCookie(*gin.Context, string)
}
type esi struct {
	iESIRoot
	base            string
	callback        string
	scope           string
	state           string
	client          string
	postDestination string
	creds           string
	link            string
}

func (x *esi) register(r *gin.RouterGroup) {
	r.GET("/", x.login)
}
func (x esi) new() *esi {

	x.scope = "scope=publicData"

	x.base = "https://login.eveonline.com/v2/oauth/authorize/?response_type=code&"
	x.callback = "redirect_uri=" + url.QueryEscape("https://localhost:8080/api/eve/")

	x.state = "state=superSecretEveApiState"
	x.client = fmt.Sprintf("client_id=%v", os.Getenv("EVE_CLIENT_ID"))
	x.postDestination = "https://login.eveonline.com/v2/oauth/token"
	x.link = fmt.Sprintf(x.base+"%v&%v&%v&%v", x.callback, x.client, x.scope, x.state)
	return &x
}
func (x *esi) getLink() string {
	return x.link
}

func (x *esi) login(c *gin.Context) {
	code := c.Query("code")
	x.creds = "grant_type=authorization_code&code=" + code
	data := fmt.Sprintf("%v:%v", os.Getenv("EVE_CLIENT_ID"), os.Getenv("EVE_SECRET_KEY"))
	encData := b64.URLEncoding.EncodeToString([]byte(data))
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", x.postDestination, strings.NewReader(x.creds))
	if err != nil {
		response{
			false,
			[]error{err, errors.New("82 from http.NewRequest")},
		}.send(c, 500)
		return
	}
	req.Header.Set("Authorization", "Basic "+encData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "login.eveonline.com")

	fmt.Println("--- making the request ---")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("--- request failes ---")
		response{
			false,
			[]string{err.Error(), "91 from client.Do"},
		}.send(c, 500)
		return
	}
	defer res.Body.Close()

	fmt.Println("--- decoding the res.Body ---")
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("--- decoding res.Body failed ---")
		response{
			false,
			[]error{err, errors.New("line 98 from ioutil.Readall")},
		}.send(c, 500)
		return
	}

	fmt.Println("--- making a new pilot ---")
	pilot, err := model.Pilot{}.Create(body)
	if err != nil {
		fmt.Println("--- making the new pilot failed ---")
		response{
			false,
			[]error{err, errors.New("line 104 trying to creat the pilot")},
		}.send(c, 500)
		return
	}

	fmt.Println("--- going to try to create a token ---")
	// create a token
	tokString, err := x.CreateToken(pilot.ID(), time.Hour*42)
	if err != nil {
		fmt.Println("--- create token failed ---")
		response{false, err}.send(c, 500)
		return
	}

	fmt.Println("--- setting cookie ---")
	// put the token in a http only cookie
	x.SetCookie(c, tokString)
}
