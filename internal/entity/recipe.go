package entity

import (
	"errors"
	"time"
)

/*
to keep it simple, the instruction will just be string to comply for the basic requirements
for the real world application might separate each of ingredient, instruction, testimonial
tab to different database table
*/
type Recipe struct {
	Id          int64
	Title       string
	Description string
	Instruction string
	Publish     *bool
	CreatedAt   time.Time
}

type RecipeDTO struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Instruction string `json:"instruction"`
	Publish     *bool  `json:"publish,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

func (r RecipeDTO) InsertValidate() error {
	if r.Title == "" {
		return errors.New("title must not be empty")
	}
	return nil
}

func (r RecipeDTO) UpdateValidate() error {
	if r.Id == 0 {
		return errors.New("id must not be empty")
	}
	if r.Title == "" {
		return errors.New("title must not be empty")
	}
	return nil
}

func (r *RecipeDTO) SetId(id int64) {
	r.Id = id
}
