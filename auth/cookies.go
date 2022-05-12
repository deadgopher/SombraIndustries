package auth

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

type CookieMaker interface {
	SetCookie(c *gin.Context, str string)
	MakeCookie(str string) (*Cookie, error)
}

type Cookie struct {
	Name     string        `json:"name"`
	Value    string        `json:"value"`
	MaxAge   time.Duration `json:"maxAge"`
	Path     string        `json:"path"`
	Domain   string        `json:"domain"`
	Secure   bool          `json:"secure"`
	HTTPOnly bool          `json:"httpOnly"`
}

func (x *JWTMaker) MakeCookie(str string) (*Cookie, error) {
	var c *Cookie
	err := json.Unmarshal([]byte(str), &c)
	return c, err
}
func (x *JWTMaker) SetCookie(c *gin.Context, str string) {
	c.SetCookie("mydamncookie", str, 100, "/", "localhost", false, true)
}
