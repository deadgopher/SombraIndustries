package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userRoot
}

func (x *userController) create(c *gin.Context) {

	user, err := x.NewUser(c)
	if err != nil {
		respond(c, 400, false, []string{err.Error(), "error in the bind user method"})
	}
	if valid, errs := user.Valid(); !valid {
		respond(c, 422, false, errs)
		return
	}
	err = user.Save()
	if err != nil {
		fmt.Println("=== ERROR TRYING TO SAVE USER :32 ===", err.Error())
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			respond(c, 422, false, foo(err))
			return
		}
		respond(c, 500, false, foo(err))
		return
	}

	respond(c, 201, true, nil)
}

func (x *userController) getOne(c *gin.Context) {
	user, err := x.NewUser(oid(c.Param("id")))
	if err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	respond(c, 200, true, user)
}

func (x *userController) getByName(c *gin.Context) {
	user, err := x.NewUser(c.Param("name"))
	if err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	respond(c, 200, true, user)
}

func (x *userController) getAll(c *gin.Context) {
	shepard, _ := x.NewUser(nil)
	users, err := shepard.GetAll()
	if err != nil {
		respond(c, 500, false, []string{err.Error(), "could not get the users"})
		return
	}
	respond(c, 200, true, users)
}

func (x *userController) delete(c *gin.Context) {
	user, _ := x.NewUser(nil)
	user.SetID(oid(c.Param("id")))
	user.Kill()
}

func (x *userController) verify(c *gin.Context) {

	user, err := x.NewUser(c)
	if err != nil {
		respond(c, 400, false, foo(err))
		return
	}

	if err := user.Verify(); err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	respond(c, 200, true, nil)

}

func (x *userController) deleteAll(c *gin.Context) {
	reaper, _ := x.NewUser(nil)
	if err := reaper.KillAll(); err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	respond(c, 200, true, nil)
}

func (x *userController) verifyToken(c *gin.Context) {
	token, err := x.VerifyToken(c.GetHeader("Authorization"))
	if err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	respond(c, 200, true, token)
}

func newUserController(r userRoot) *userController {
	return &userController{r}
}

// Login : Log in
func (x *userController) login(c *gin.Context) {
	user, err := x.NewUser(c)
	if err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	if err := user.Login(); err != nil {
		respond(c, 400, false, []string{err.Error(), ""})
	}
	ts, err := x.CreateToken(user, time.Hour*24)
	if err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	respond(c, 200, true, ts)
}

func (x *userController) registerRoutes(r *gin.RouterGroup) {
	r.POST("/", x.create)
	r.GET("/:id", x.getOne)
	r.GET("/n/:name", x.getByName)
	r.GET("/", x.getAll)
	r.DELETE("/:id", x.protect(x.delete))
	r.DELETE("/", x.deleteAll)
	r.POST("/login", x.login)
	r.GET("/auth", x.verifyToken)
	r.PUT("/verify", x.protect(x.verify))
}
