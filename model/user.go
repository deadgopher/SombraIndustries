package model

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type IUser interface {
	SetID(primitive.ObjectID)
	Valid() (bool, []string)
	Save() error
	Login() error
	GetAll() ([]*User, error)
	Kill() error
	KillAll() error
	Verify() error
}

type User struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	Name            string `json:"name" bson:"name"`
	Password        string `json:"password" bson:"password"`
	Email           string `json:"email" bson:"email"`
	ConfirmPassword string `json:"cpassword,omitempty" bson:"-"`
	Verified        bool   `json:"verified" bson:"verified"`

	col *mongo.Collection `json:"-" bson:"-"`
	ctx context.Context   `json:"-" bson:"-"`

	CreatedAt time.Time `json:"_createdAt,omitempty" bson:"_createdAt,omitempty"`
	UpdatedAt time.Time `json:"_updatedAt,omitempty" bson:"_updatedAt,omitempty"`
}

func (x *User) SetID(id primitive.ObjectID) {
	x.ID = id
}

func (x *User) Valid() (bool, []string) {
	var verrors []string
	if len(x.Name) < 3 {
		verrors = append(verrors, "username too short")
	}
	if err := x.passwordIsValid(); err != nil {
		verrors = append(verrors, err.Error())
	}

	if len(verrors) != 0 {
		return false, verrors
	}
	return true, nil
}

func (x *User) passwordIsValid() error {
	if len(x.Password) < 8 {
		return errors.New("username too short")
	}
	if x.Password != x.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	return nil

}

func (x *User) hashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(x.Password), 14)
	if err != nil {
		return err
	}
	x.Password = string(bytes)
	return nil
}

func (x *User) passwordMatches(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(x.Password), []byte(password))
	return err == nil
}

func (x *User) Save() error {
	if err := x.isUnique(); err != nil {
		return err
	}
	x.CreatedAt = time.Now()
	x.UpdatedAt = time.Now()
	x.Verified = false
	// encrypt password
	if err := x.hashPassword(); err != nil {
		return err
	}
	_, err := x.col.InsertOne(x.ctx, x)
	return err
}

func (x *User) isUnique() error {
	filter1 := bson.M{
		"name": x.Name,
	}
	filter2 := bson.M{
		"email": x.Email,
	}
	if err := x.col.FindOne(x.ctx, filter1).Decode(&x); err == nil {
		return errors.New("username already exists")
	}
	if err := x.col.FindOne(x.ctx, filter2).Decode(&x); err == nil {
		return errors.New("email already exists")
	}
	return nil
}

// Login :
func (x *User) Login() error {
	var tmpUser *User
	filter := bson.M{
		"name": x.Name,
	}
	if err := x.col.FindOne(x.ctx, filter).Decode(&tmpUser); err != nil {
		return errors.New("could not find that User")
	}
	if !x.passwordMatches(x.Password) {
		return errors.New("incorrect password")
	}
	return nil

}
func (x *User) GetAll() ([]*User, error) {
	var users []*User
	cur, err := x.col.Find(x.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(x.ctx)
	for cur.Next(x.ctx) {
		var user *User
		if err := cur.Decode(&user); err != nil {
			return nil, err
		}
		user.bindToDB(x.col, x.ctx)
		users = append(users, user)
	}
	return users, nil
}

func (x *User) Kill() error {
	_, err := x.col.DeleteOne(x.ctx, bson.M{"_id": x.ID})
	return err
}
func (x *User) KillAll() error {
	_, err := x.col.DeleteMany(x.ctx, bson.M{})
	return err
}
func (x *User) Verify() error {
	update := bson.M{
		"$set": bson.M{
			"verified": true,
		},
	}
	_, err := x.col.UpdateByID(x.ctx, x.ID, update)
	return err
}

func (x *User) bindToDB(col *mongo.Collection, ctx context.Context) {
	x.col = col
	x.ctx = ctx
}
