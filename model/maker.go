package model

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ModelMaker interface {
	NewUser(interface{}) (*User, error)
	NewComment(interface{}) (*Comment, error)
	NewPost(c interface{}) (*Post, error)
}

type maker struct {
	*mongo.Database
	ctx context.Context
}

func NewMaker(db *mongo.Database, ctx context.Context) *maker {
	x := &maker{db, ctx}
	return x
}

func (x *maker) NewUser(c interface{}) (*User, error) {
	var user *User
	user.col = x.Collection("users")
	user.ctx = x.ctx
	if _, ok := c.(primitive.ObjectID); ok {
		if err := x.Collection("users").FindOne(x.ctx, bson.M{"_id": c.(primitive.ObjectID)}).Decode(&user); err != nil {
			return nil, err
		}
	}
	if _, ok := c.(*gin.Context); ok {
		if err := c.(*gin.Context).BindJSON(&user); err != nil {
			return nil, err
		}
	}
	if _, ok := c.(string); ok {
		if err := x.Collection("users").FindOne(x.ctx, bson.M{"name": c.(string)}).Decode(&user); err != nil {
			return nil, err
		}
	}
	return user, nil

}
func (x *maker) NewComment(c interface{}) (*Comment, error) {
	var comment *Comment
	comment.col = x.Collection("comments")
	comment.ctx = x.ctx
	if _, ok := c.(primitive.ObjectID); ok {
		if err := x.Collection("comments").FindOne(x.ctx, bson.M{"_id": c.(primitive.ObjectID)}).Decode(&comment); err != nil {
			return nil, err
		}
	}
	if _, ok := c.(*gin.Context); ok {
		if err := c.(*gin.Context).BindJSON(&comment); err != nil {
			return nil, err
		}
	}
	return comment, nil
}

func (x *maker) NewPost(c interface{}) (*Post, error) {
	var post *Post
	post.col = x.Collection("posts")
	post.ctx = x.ctx
	if _, ok := c.(*gin.Context); ok {
		if err := c.((*gin.Context)).BindJSON(&post); err != nil {
			return nil, err
		}
	}
	if _, ok := c.(primitive.ObjectID); ok {
		if err := x.Collection("posts").FindOne(x.ctx, bson.M{"_id": c.(primitive.ObjectID)}).Decode(&post); err != nil {
			return nil, err
		}
	}
	return post, nil
}
