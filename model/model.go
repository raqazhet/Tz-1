package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	ActiveAt  time.Time          `json:"activeAt" bson:"activeAt"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	Status    string             `json:"status" bson:"status"`
}
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}
type SuccessResponse struct {
	Message string `json:"message"`
}
