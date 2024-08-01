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

type GetAllAuthor struct {
	CurrentPage *int      `json:"current_page"`
	NextPage    *int      `json:"next_page"`
	Author      *[]Auther `json:"auther"`
}

type AutherUpdate struct {
	FirstName *string   `json:"first_name" validate:"required"`
	LastName  *string   `json:"last_name" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    *string   `json:"user_id"`
}

type Book struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      *string            `json:"name" validate:"required"`
	ISBN      *string            `json:"isbn" validate:"required"`
	AuthorId  *string            `json:"author_id" validate:"required"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	BookId    *string            `json:"book_id"`
}

type BookResponse struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      *string            `json:"name" validate:"required"`
	ISBN      *string            `json:"isbn" validate:"required"`
	Auther    *Auther            `json:"author"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	BookId    *string            `json:"book_id"`
}

type GetBook struct {
	Message *string       `json:"message"`
	Success *bool         `json:"success"`
	Book    *BookResponse `json:"book"`
}
