package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"germ/auth"
	"germ/model"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	Register(*gin.RouterGroup)
}
type iController interface {
	register(*gin.RouterGroup)
}
type iRoot interface {
}

type root struct {
	iESIController
	iPilotController
	iPostController
	iCommentController
	auth.TokenMaker
}

func New(db *sql.DB) Controller {
	model.Init(db)
	x := &root{}
	x.TokenMaker = auth.NewJWTMaker(os.Getenv("SECRET_KEY"))

	x.iESIController = esi{iESIRoot: x}.new()
	x.iPilotController = &pilots{iPilotRoot: x}
	x.iPostController = &posts{iPostRoot: x}
	x.iCommentController = &comments{iCommentRoot: x}

	return x

}

func (x *root) Register(r *gin.RouterGroup) {
	x.iESIController.register(r.Group("eve"))
	x.iPilotController.register(r.Group("pilots"))
	x.iPostController.register(r.Group("posts"))
	x.iCommentController.register(r.Group("comments"))

	r.GET("/", x.handshake)
}

func (x *root) handshake(c *gin.Context) {

	fmt.Println("trying to handshake...")

	fmt.Println("getting cookie")
	cookieString, err := c.Cookie("mydamncookie")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("---there was no cookie " + err.Error())
			response{true, struct {
				Link  string       `json:"link"`
				Pilot *model.Pilot `json:"pilot"`
			}{
				x.getLink(),
				nil,
			}}.send(c, http.StatusOK)
			return
		}
		fmt.Println("---something else went wrong " + err.Error())
		response{false, err.Error()}.send(c, http.StatusInternalServerError)
		return

	}
	cookie, err := x.MakeCookie(cookieString)
	if err != nil {
		fmt.Println("---error unmarshaling the cookie " + err.Error())
		response{false, err}.send(c, http.StatusInternalServerError)
		return
	}
	payload, err := x.VerifyToken(cookie.Value)
	if err != nil {
		fmt.Println("---error verifying the token string inside the cookie else went wrong " + err.Error())
		response{false, err}.send(c, http.StatusInternalServerError)
		return
	}

	fmt.Println("---making res data ")

	fmt.Println("---responding success")
	response{true, struct {
		Link  string `json:"link"`
		Pilot string `json:"pilot"`
	}{
		x.getLink(),
		payload.Data,
	}}.send(c, http.StatusOK)
}

// func (x *rootController) protect(next gin.HandlerFunc) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokString := c.GetHeader("Authorization")
// 		payload, err := x.VerifyToken(tokString)
// 		if err != nil {
// 			respond(c, 400, false, []string{err.Error(), "not authorized!"})
// 			return
// 		}
// 		c.Header("Authorization", payload.ID.String())
// 		next(c)
// 	}
// }
