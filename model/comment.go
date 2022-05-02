package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IComment interface {
	Save() error
	GetAll(primitive.ObjectID) ([]*Comment, error)
	Valid() (bool, []string)
	Update() error
	Kill() error
	KillAll() error
}

type Comment struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	Author   string             `json:"author" bson:"author"`
	AuthorID primitive.ObjectID `json:"_author" bson:"_author"`
	ParentID primitive.ObjectID `json:"_parent" bson:"_parent"`
	Body     string             `json:"body" bson:"body"`

	col *mongo.Collection `json:"-" bson:"-"`
	ctx context.Context   `json:"-" bson:"-"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (x *Comment) Save() error {

	x.CreatedAt = time.Now()
	x.UpdatedAt = time.Now()

	res, err := x.col.InsertOne(x.ctx, x)
	if err != nil {
		return err
	}
	x.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (x *Comment) GetAll(id primitive.ObjectID) ([]*Comment, error) {
	var comments []*Comment
	cur, err := x.col.Find(x.ctx, bson.M{"_parent": id})
	if err != nil {
		return nil, err
	}
	defer cur.Close(x.ctx)
	for cur.Next(x.ctx) {
		var Comment *Comment
		if err := cur.Decode(&Comment); err != nil {
			return nil, err
		}
		Comment.bindToDB(x.col, x.ctx)
		comments = append(comments, Comment)
	}
	return comments, nil

}

func (x *Comment) Valid() (bool, []string) {
	var verrors []string

	if len(x.Body) == 0 {
		verrors = append(verrors, "you forgot to write something fuckface")
	}

	if len(verrors) != 0 {
		return false, verrors
	}
	return true, nil
}

func (x *Comment) Update() error {
	update := bson.M{
		"$set": bson.M{
			"body": x.Body,
		},
	}
	_, err := x.col.UpdateByID(x.ctx, x.ID, update)
	return err
}

func (x *Comment) Kill() error {
	_, err := x.col.DeleteOne(x.ctx, bson.M{"_id": x.ID})
	return err
}
func (x *Comment) KillAll() error {
	_, err := x.col.DeleteMany(x.ctx, bson.M{})
	return err
}

func (x *Comment) bindToDB(col *mongo.Collection, ctx context.Context) {
	x.col = col
	x.ctx = ctx
}
