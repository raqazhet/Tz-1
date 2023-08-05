package model

import "time"

type Todo struct {
	ID       string    `json:"id,omitempty" bson:"_id,omitempty"`
	Title    string    `json:"title" validate:"required,max=200"`
	ActiveAt time.Time `json:"activeAt" validate:"required"`
	Status   string    `json:"status"`
}
