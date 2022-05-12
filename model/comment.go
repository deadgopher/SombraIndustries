package model

import (
	"errors"
	"strconv"
	"time"
)

type Comment struct {
	id int

	author int
	parent int
	body   string

	createdAt time.Time
	updatedAt time.Time
}

func (x Comment) Create(i interface{}) (*Comment, error) {

	// Get One
	if id, ok := i.(string); ok {
		if intID, err := strconv.Atoi(id); err != nil {
			return nil, err
		} else {
			x.id = intID
		}
		rows, err := eveHQ.Query(`
		SELECT *
		FROM comments
		WHERE id = %v
		`, x.id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		if err := rows.Scan(&x); err != nil {
			return nil, err
		}
	}

	return &x, nil
}

func (x Comment) Read(id string) ([]*Comment, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	var comments []*Comment
	rows, err := eveHQ.Query(`
	SELECT *
	FROM comments
	WHERE parent = %v
	`, intID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment *Comment
		if err := rows.Scan(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (x *Comment) Update() error {
	_, err := eveHQ.Query(`
	UPDATE  comments
	SET body = %v
	WHERE id = %v
	`, x.body, x.id)
	return err
}

func (x *Comment) Destroy() error {
	_, err := eveHQ.Query(`
	DELETE *
	FROM comments
	WHERE id = %v
	`, x.id)
	return err
}

func (x Comment) Purge(id string) error {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	_, err = eveHQ.Query(`
	DELETE *
	FROM comments
	WHERE parent = %v
	`, intID)
	return err
}

func (x *Comment) Save() error {
	_, err := eveHQ.Query(`
	INSERT INTO comments
	VALUES(%v,%v,%v,%v,%v,%v)
	`, x)
	return err
}

func (x *Comment) Validate() []error {
	var errs []error

	if len(x.body) == 0 {
		errs = append(errs, errors.New("you forgot to write something fuckface"))
	}

	return errs
}
