package controller

import (
	"github.com/gin-gonic/gin"
)

type postController struct {
	root
}

// Create
func (x *postController) create(c *gin.Context) {
	post, err := x.NewPost(c)
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

func (x *postController) getOne(c *gin.Context) {
	post, err := x.NewPost(oid(c.Param("id")))
	if err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	respond(c, 200, true, post)
}
func (x *postController) getAll(c *gin.Context) {
	shepard, _ := x.NewPost(nil)
	posts, err := shepard.GetAll()
	if err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	respond(c, 200, true, posts)
}

func (x *postController) update(c *gin.Context) {

}
func (x *postController) delete(c *gin.Context) {
	post, err := x.NewPost(oid(c.Param("id")))
	if err != nil {
		respond(c, 400, false, foo(err))
		return
	}
	if err := post.Kill(); err != nil {
		respond(c, 500, false, foo(err))
		return
	}
	respond(c, 200, true, nil)
}

func (x *postController) deleteAll(c *gin.Context) {
	reaper, _ := x.NewPost(nil)
	if err := reaper.KillAll(); err != nil {
		respond(c, 400, false, []string{err.Error(), "dumb ass"})
		return
	}
	respond(c, 200, true, "good job!")
}

func (x *postController) registerRoutes(r *gin.RouterGroup) {
	r.POST("/", x.protect(x.create))
	r.GET("/:id", x.getOne)
	r.GET("/", x.getAll)
	r.PUT("/", x.protect(x.update))
	r.DELETE("/:id", x.protect(x.delete))
	r.DELETE("/", x.deleteAll)
}

func newPostController(r root) *postController {
	x := &postController{
		r,
	}
	return x
}
