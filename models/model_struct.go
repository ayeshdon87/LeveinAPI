package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auther struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName *string            `json:"first_name" validate:"required"`
	LastName  *string            `json:"last_name" validate:"required"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	UserId    *string            `json:"user_id"`
}

type AutherCreateSuccess struct {
	Id      *string `json:"id"`
	Message *string `json:"message"`
	Success *bool   `json:"success"`
}

type GetAuthor struct {
	Message *string `json:"message"`
	Success *bool   `json:"success"`
	Author  *Auther `json:"auther"`
}
