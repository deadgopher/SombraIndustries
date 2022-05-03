package controller

import (
	"github.com/gin-gonic/gin"
)

type iPosts interface {
	RegisterPosts(r *gin.RouterGroup)
}

type posts struct {
	postRoot
}

func newPostController(r postRoot) *posts {
	x := &posts{
		r,
	}
	return x
}

// Create
func (x *posts) create(c *gin.Context) {
	post, err := x.Post(c)
	if err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	if validationErrors := post.Valid(); validationErrors != nil {
		respond(c, 400, false, validationErrors)
		return
	}
	if err := post.Save(); err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	respond(c, 201, true, post)
}

func (x *posts) getOne(c *gin.Context) {
	post, err := x.Post(oid(c.Param("id")))
	if err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	respond(c, 200, true, post)
}
func (x *posts) getAll(c *gin.Context) {
	shepard, _ := x.Post(nil)
	posts, err := shepard.Get()
	if err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	respond(c, 200, true, posts)
}

func (x *posts) update(c *gin.Context) {

}
func (x *posts) delete(c *gin.Context) {
	post, err := x.Post(oid(c.Param("id")))
	if err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	if err := post.Delete(); err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	respond(c, 200, true, nil)
}

func (x *posts) deleteAll(c *gin.Context) {
	if err := x.DeletePosts(); err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	respond(c, 200, true, nil)
}

func (x *posts) RegisterPosts(r *gin.RouterGroup) {
	r.POST("/", x.create)
	r.GET("/:id", x.getOne)
	r.GET("/", x.getAll)
	r.PUT("/", x.update)
	r.DELETE("/:id", x.delete)
	r.DELETE("/", x.deleteAll)
}
