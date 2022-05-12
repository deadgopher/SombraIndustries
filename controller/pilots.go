package controller

import (
	"germ/model"

	"github.com/gin-gonic/gin"
)

// sharing
type iPilotController interface {
	iController
}

// requires
type iPilotRoot interface {
}
type pilots struct {
	iPilotRoot
}

func (x *pilots) register(r *gin.RouterGroup) {
	r.GET("/:id", x.getOne)
	r.GET("/n/:name", x.getByName)
	r.GET("/", x.getAll)
	r.DELETE("/:id", x.delete)
	r.DELETE("/", x.deleteAll)
}

func (x *pilots) getOne(c *gin.Context) {
	pilot, err := model.Pilot{}.Create(c.Param("id"))
	if err != nil {
		response{false, []string{err.Error()}}.send(c, 400)
		return
	}
	response{true, pilot}.send(c, 200)
}

func (x *pilots) getByName(c *gin.Context) {
	pilot, err := model.Pilot{}.Create(c.Param("id"))
	if err != nil {
		response{false, []string{err.Error()}}.send(c, 400)
		return
	}
	response{true, pilot}.send(c, 200)
}

func (x *pilots) getAll(c *gin.Context) {
	pilots, err := model.Pilot{}.Read()
	if err != nil {
		response{false, []string{err.Error()}}.send(c, 500)
		return
	}
	response{true, pilots}.send(c, 200)
}

func (x *pilots) delete(c *gin.Context) {
	pilot, err := model.Pilot{}.Create(c.Param("id"))
	if err != nil {
		response{false, []string{err.Error()}}.send(c, 400)
		return
	}
	if err := pilot.Destroy(); err != nil {
		response{false, []string{err.Error()}}.send(c, 500)
		return
	}
	response{true, nil}.send(c, 200)
}

func (x *pilots) deleteAll(c *gin.Context) {

	err := model.Pilot{}.Purge()
	if err != nil {
		response{false, []string{err.Error()}}.send(c, 500)
		return
	}
	response{true, nil}.send(c, 200)
}
