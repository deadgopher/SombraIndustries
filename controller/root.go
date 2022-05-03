package controller

import (
	"context"
	"os"
	"time"

	"germ/model"
	"germ/token"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func foo(bar error) []string {
	return []string{bar.Error()}
}

func oid(id string) primitive.ObjectID {
	x, _ := primitive.ObjectIDFromHex(id)
	return x
}

func respond(c *gin.Context, code int, success bool, payload interface{}) {
	r := gin.H{
		"success": success,
		"payload": payload,
	}
	c.JSON(code, r)
}

type Routes interface {
	RegisterPilots(*gin.RouterGroup)
	RegisterPosts(*gin.RouterGroup)
	RegisterComments(*gin.RouterGroup)
	RegisterViews(*gin.RouterGroup)
	RegisterEve(*gin.RouterGroup)
}

type rootController struct {
	model.ModelMaker
	iEve
	iPilots
	iPosts
	iComments
	iViews
	token.TokenMaker
}

func New(db *mongo.Database, ctx context.Context) Routes {
	x := &rootController{}

	x.ModelMaker = model.NewMaker(db, ctx)
	x.iPilots = newPilotController(x)
	x.iPosts = newPostController(x)
	x.iComments = newCommentController(x)
	x.iViews = newViewController(x)
	x.TokenMaker = token.NewJWTMaker(os.Getenv("SECRET_KEY"))
	x.iEve = newEveController(x)
	return x
}

type root interface {
	// protect(gin.HandlerFunc) gin.HandlerFunc
	setData(interface{})
	PilotFromToken([]byte) (*model.Pilot, error)
	Comment(interface{}) (*model.Comment, error)
	Post(interface{}) (*model.Post, error)
	MakeRequest() ([]byte, error)
	CreateToken(data *model.Pilot, exp time.Duration) (string, error)
	show(gin.ResponseWriter, string)
}
type userRoot interface {
	root
	VerifyToken(string) (*token.Payload, error)
	DeleteUsers() error
}
type postRoot interface {
	root
	DeletePosts() error
}
type commentRoot interface {
	root
	DeleteComments() error
}
type viewRoot interface {
	root
	esiLink() string
}
type eveRoot interface {
	root
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
