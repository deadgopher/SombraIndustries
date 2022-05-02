package model

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IPost interface {
	Save() error
	GetAll() ([]*Post, error)
	Update() error
	Kill() error
	KillAll() error
	Valid() []error
}

type Post struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	Author   string             `json:"author" bson:"author"`
	AuthorID primitive.ObjectID `json:"_author" bson:"_author"`
	Title    string             `json:"title" bson:"title"`
	Body     string             `json:"body" bson:"body"`
	Comments []*Comment         `json:"comments" bson:"-"`

	col *mongo.Collection `json:"-" bson:"-"`
	ctx context.Context   `json:"-" bson:"-"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (x *Post) Save() error {
	x.CreatedAt = time.Now()
	x.UpdatedAt = time.Now()
	res, err := x.col.InsertOne(x.ctx, x)
	if err != nil {
		return err
	}
	x.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (x *Post) GetAll() ([]*Post, error) {
	var posts []*Post
	cur, err := x.col.Find(x.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(x.ctx)
	for cur.Next(x.ctx) {
		var Post *Post
		if err := cur.Decode(&Post); err != nil {
			return nil, err
		}
		Post.bindToDB(x.col, x.ctx)
		posts = append(posts, Post)
	}
	return posts, nil
}

func (x *Post) Update() error {
	update := bson.M{
		"$set": bson.M{
			"title": x.Title,
			"body":  x.Body,
		},
	}
	_, err := x.col.UpdateByID(x.ctx, x.ID, update)
	return err
}

func (x *Post) Kill() error {
	_, err := x.col.DeleteOne(x.ctx, bson.M{"_id": x.ID})
	return err
}
func (x *Post) KillAll() error {
	_, err := x.col.DeleteMany(x.ctx, bson.M{})
	return err
}

func (x *Post) Valid() []error {
	var errs []error
	if len(x.Body) < 3 {
		errs = append(errs, errors.New("this Post is too damn short"))
	}

	if len(errs) != 0 {
		return errs
	}
	return nil
}

func (x *Post) bindToDB(col *mongo.Collection, ctx context.Context) {
	x.col = col
	x.ctx = ctx
}
