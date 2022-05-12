package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type Post struct {
	id     int
	author int
	title  string
	body   string

	createdAt time.Time
	updatedAt time.Time
}

func (x *Post) Save() error {
	x.createdAt = time.Now()
	x.updatedAt = time.Now()
	q := fmt.Sprintf(`
	INSERT INTO posts
	VALUES(%v,%v,%v,%v,%v)
	`, x.author, x.title, x.body, x.createdAt, x.updatedAt)
	insert, err := eveHQ.Query(q)
	if err != nil {
		return err
	}
	defer insert.Close()

	err = insert.Scan(&x.id)
	return err
}

func (x Post) Create(i interface{}) (*Post, error) {
	if _, ok := i.(string); ok {
		bytes := []byte(i.(string))
		if err := json.Unmarshal(bytes, &x.id); err != nil {
			return nil, err
		}
	}
	if _, ok := i.(*gin.Context); ok {
		c := i.(*gin.Context)
		if err := c.BindJSON(&x); err != nil {
			return nil, err
		}
	}

	return &x, nil
}

func (x Post) Read() ([]*Post, error) {
	var posts []*Post

	rows, err := eveHQ.Query(`
	SELECT *
	FROM posts
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post *Post
		if err := rows.Scan(&post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (x *Post) Update() error {
	_, err := eveHQ.Query(`
	UPDATE posts
	SET title = %v, body = %v, updatedAt = %v
	WHERE id = %v
	`, x.title, x.body, x.updatedAt, x.id)
	return err
}

func (x *Post) Destroy() error {
	_, err := eveHQ.Query(`
	DELETE * 
	FROM posts
	WHERE id = %v
	`, x.id)
	return err
}
func (x Post) Purge() error {
	return nil
}

func (x *Post) Validate() []string {
	var errs []string
	if len(x.body) < 3 {
		errs = append(errs, "this post is too damn short")
	}

	if len(errs) != 0 {
		return errs
	}
	return nil
}
