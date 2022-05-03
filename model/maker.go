package model

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _database *mongo.Database = 
var _root_context context.Context

type ModelMaker interface {
	PilotFromToken([]byte) (*Pilot, error)
	Comment(interface{}) (*Comment, error)
	Post(c interface{}) (*Post, error)
	DeleteUsers() error
	DeletePosts() error
	DeleteComments() error
}

type maker struct {
	*mongo.Database
	ctx context.Context
}

func NewMaker(db *mongo.Database, ctx context.Context) *maker {
	x := &maker{db, ctx}
	return x
}

func (x *maker) PilotFromToken(t []byte) (*Pilot, error) {
	var pilot *Pilot
	pilot.Collection = x.Collection("users")
	pilot.ctx = x.ctx

	if err := json.Unmarshal(t, &pilot.ESIToken); err != nil {
		return nil, err
	}
	return pilot, nil
}
func (x *maker) Comment(c interface{}) (*Comment, error) {
	var comment *Comment
	comment.Collection = x.Collection("comments")
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

func (x *maker) Post(c interface{}) (*Post, error) {
	var post *Post
	post.Collection = x.Collection("posts")
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
func (x *maker) DeleteUsers() error {
	_, err := x.Collection("users").DeleteMany(x.ctx, bson.M{})
	return err
}
func (x *maker) DeletePosts() error {
	_, err := x.Collection("posts").DeleteMany(x.ctx, bson.M{})
	return err
}
func (x *maker) DeleteComments() error {
	_, err := x.Collection("comments").DeleteMany(x.ctx, bson.M{})
	return err
}
