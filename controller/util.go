package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func oid(id string) primitive.ObjectID {
	x, _ := primitive.ObjectIDFromHex(id)
	return x
}

type response struct {
	Success bool        `json:"success"`
	Payload interface{} `json:"payload"`
}

func (x response) send(c *gin.Context, code int) {
	c.JSON(code, x)
}
