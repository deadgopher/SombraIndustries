package controller

import (
	"github.com/gin-gonic/gin"
)

type iPilots interface {
	RegisterPilots(r *gin.RouterGroup)
}

type pilots struct {
	userRoot
}

func newPilotController(r userRoot) *pilots {
	return &pilots{r}
}
func (x *pilots) RegisterPilots(r *gin.RouterGroup) {
	// r.POST("/", x.create)
	r.GET("/:id", x.getOne)
	r.GET("/n/:name", x.getByName)
	r.GET("/", x.getAll)
	r.DELETE("/:id", x.delete)
	r.DELETE("/", x.deleteAll)
	// r.POST("/login", x.login)
	r.GET("/auth", x.verifyToken)
}

func (x *pilots) getOne(c *gin.Context) {

}

func (x *pilots) getByName(c *gin.Context) {

}

func (x *pilots) getAll(c *gin.Context) {

}

func (x *pilots) delete(c *gin.Context) {

}

func (x *pilots) deleteAll(c *gin.Context) {

}

func (x *pilots) verifyToken(c *gin.Context) {
	token, err := x.VerifyToken(c.GetHeader("Authorization"))
	if err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	respond(c, 200, true, token)
}

// Login : Log in
func (x *pilots) login(c *gin.Context) {

}
