package controller

import (
	"germ/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// sharing
type iCommentController interface {
	iController
}

// requires
type iCommentRoot interface {
}
type comments struct {
	iCommentRoot
}

func (x *comments) register(r *gin.RouterGroup) {
	r.POST("/", x.create)
	r.GET("/:id", x.getAll)
	r.PUT("/", x.update)
	r.DELETE("/:id", x.delete)
	r.DELETE(("/all/:id"), x.deleteAll)
}

func (x *comments) create(c *gin.Context) {
	comment, err := model.Comment{}.Create(c)
	if err != nil {
		response{false, []string{err.Error()}}.send(c, http.StatusBadRequest)
		return
	}
	if errs := comment.Validate(); errs != nil {
		response{false, []string{err.Error()}}.send(c, http.StatusBadRequest)
		return
	}
	if err := comment.Save(); err != nil {
		response{false, []string{err.Error()}}.send(c, http.StatusInternalServerError)
	}
	response{true, comment}.send(c, http.StatusCreated)
}

func (x *comments) getAll(c *gin.Context) {

	comments, err := model.Comment{}.Read(c.Param("id"))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			response{false, "no comment exists with that id"}.send(c, http.StatusNoContent)
			return
		} else {
			response{false, []string{err.Error()}}.send(c, http.StatusInternalServerError)
			return
		}
	}
	response{true, comments}.send(c, http.StatusFound)
}

func (x *comments) update(c *gin.Context) {

	comment, err := model.Comment{}.Create(c)
	if err != nil {
		response{false, []string{err.Error()}}.send(c, http.StatusBadRequest)
		return
	}
	if err := comment.Update(); err != nil {
		response{false, []string{err.Error()}}.send(c, http.StatusInternalServerError)
		return
	}
	response{true, comment}.send(c, http.StatusOK)
}
func (x *comments) delete(c *gin.Context) {
	comment, err := model.Comment{}.Create(c.Param("id"))
	if err != nil {
		response{false, []string{err.Error()}}.send(c, http.StatusBadRequest)
		return
	}
	if err := comment.Destroy(); err != nil {
		response{false, []string{err.Error()}}.send(c, http.StatusInternalServerError)
		return
	}
	response{true, nil}.send(c, http.StatusOK)
}
func (x *comments) deleteAll(c *gin.Context) {
	err := model.Comment{}.Purge(c.Param("id"))
	if err != nil {
		response{false, []string{err.Error()}}.send(c, http.StatusBadRequest)
		return
	}
	response{true, nil}.send(c, http.StatusOK)
}
