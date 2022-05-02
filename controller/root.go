package controller

import (
	"context"
	"fmt"
	"os"
	"time"

	"germ/model"
	"germ/token"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type rootController struct {
	model.ModelMaker
	iEveController
	users    controller
	posts    controller
	comments controller
	views    controller
	token.TokenMaker
}

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
	RegisterUsers(*gin.RouterGroup)
	RegisterPosts(*gin.RouterGroup)
	RegisterComments(*gin.RouterGroup)
	RegisterViews(*gin.RouterGroup)
	RegisterEve(*gin.RouterGroup)
}

type root interface {
	protect(gin.HandlerFunc) gin.HandlerFunc
	NewUser(interface{}) (*model.User, error)
	NewComment(interface{}) (*model.Comment, error)
	NewPost(interface{}) (*model.Post, error)
}
type userRoot interface {
	root
	VerifyToken(string) (*token.Payload, error)
	CreateToken(data *model.User, exp time.Duration) (string, error)
}
type viewRoot interface {
	root
	Link() string
}

type controller interface {
	registerRoutes(*gin.RouterGroup)
}

type iEveController interface {
	controller
	Link() string
}

func (x *rootController) protect(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokString := c.GetHeader("Authorization")
		payload, err := x.VerifyToken(tokString)
		if err != nil {
			respond(c, 400, false, []string{err.Error(), "not authorized!"})
			return
		}
		c.Header("Authorization", payload.ID.String())
		next(c)
	}
}

func New(db *mongo.Database, ctx context.Context) Routes {
	x := &rootController{}

	x.ModelMaker = model.NewMaker(db, ctx)
	fmt.Println("made model maker")
	x.users = newUserController(x)
	fmt.Println("made user controller")
	x.posts = newPostController(x)
	fmt.Println("made post controller")
	x.comments = newCommentController(x)
	fmt.Println("made comment controller")
	x.views = newViewController(x)
	fmt.Println("made view controller")
	x.TokenMaker = token.NewJWTMaker(os.Getenv("SECRET_KEY"))
	fmt.Println("made token maker")
	x.iEveController = newEveController(x)
	return x
}

func (x *rootController) RegisterUsers(g *gin.RouterGroup) {
	x.users.registerRoutes(g)
}
func (x *rootController) RegisterPosts(g *gin.RouterGroup) {
	x.posts.registerRoutes(g)
}
func (x *rootController) RegisterComments(g *gin.RouterGroup) {
	x.comments.registerRoutes(g)
}
func (x *rootController) RegisterViews(g *gin.RouterGroup) {
	x.views.registerRoutes(g)
}

func (x *rootController) RegisterEve(g *gin.RouterGroup) {
	x.iEveController.registerRoutes(g)
}
