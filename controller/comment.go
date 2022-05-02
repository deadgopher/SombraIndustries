package controller

import (
	"github.com/gin-gonic/gin"
)

type commentController struct {
	root
}

func (x *commentController) create(c *gin.Context) {

}

func (x *commentController) getAll(c *gin.Context) {
	shepard, _ := x.NewComment(nil)
	comments, err := shepard.GetAll(oid(c.Param("id")))
	if err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	respond(c, 200, true, comments)
}

func (x *commentController) update(c *gin.Context) {

	comment, err := x.NewComment(c)
	if err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	if err := comment.Update(); err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	respond(c, 200, true, nil)
}
func (x *commentController) delete(c *gin.Context) {
	comment, err := x.NewComment(oid(c.Param("id")))
	if err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	if err := comment.Kill(); err != nil {
		respond(c, 500, false, foo(err))
	}
	respond(c, 200, true, nil)

}

func (x *commentController) deleteAll(c *gin.Context) {
	reaper, _ := x.NewComment(nil)
	if err := reaper.KillAll(); err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	respond(c, 200, true, "all comments gone!")
}

func (x *commentController) registerRoutes(r *gin.RouterGroup) {
	r.POST("/", x.protect(x.create))
	r.GET("/:id", x.getAll)
	r.PUT("/", x.protect(x.update))
	r.DELETE("/:id", x.protect(x.delete))
	r.DELETE("/", x.deleteAll)
}

func newCommentController(r root) *commentController {
	x := &commentController{
		r,
	}
	return x
}
