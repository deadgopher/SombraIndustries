package controller

import (
	"github.com/gin-gonic/gin"
)

type iComments interface {
	RegisterComments(r *gin.RouterGroup)
}

type comments struct {
	commentRoot
}

func newCommentController(r commentRoot) *comments {
	x := &comments{
		r,
	}
	return x
}

func (x *comments) create(c *gin.Context) {
	comment, err := x.Comment(c)
	if err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	if valid, errs := comment.Valid(); !valid {
		respond(c, 400, false, errs)
		return
	}
	if err := comment.Save(); err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	respond(c, 201, true, comment)
}

func (x *comments) getAll(c *gin.Context) {
	shepard, _ := x.Comment(nil)
	comments, err := shepard.Get(oid(c.Param("id")))
	if err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	respond(c, 200, true, comments)
}

func (x *comments) update(c *gin.Context) {

	comment, err := x.Comment(c)
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
func (x *comments) delete(c *gin.Context) {
	comment, err := x.Comment(oid(c.Param("id")))
	if err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	if err := comment.Delete(); err != nil {
		respond(c, 500, false, foo(err))
	}
	respond(c, 200, true, nil)

}

func (x *comments) deleteAll(c *gin.Context) {

}

func (x *comments) RegisterComments(r *gin.RouterGroup) {
	r.POST("/", x.create)
	r.GET("/:id", x.getAll)
	r.PUT("/", x.update)
	r.DELETE("/:id", x.delete)
	r.DELETE("/", x.deleteAll)
}
