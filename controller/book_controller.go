package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ayeshdon87/LeveinAPI/database"
	"github.com/ayeshdon87/LeveinAPI/models"
	"github.com/ayeshdon87/LeveinAPI/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var bookValidate = validator.New()

var bookCollection *mongo.Collection = database.OpentCollection(database.Client,
	utils.BOOK_TABLE)

// var authorCollection *mongo.Collection = database.OpentCollection(database.Client,
// 	utils.AUTHER_TABLE)

func AddBook() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var bookModl models.Book

		err := c.BindJSON(&bookModl)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		validationError := validate.Struct(bookModl)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.INVALID_RQUEST})
			defer cancel()
			return
		}

		//======================================

		count, err := authorCollection.CountDocuments(ctx, bson.M{"userid": bookModl.AuthorId})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{utils.ERROR: utils.AUTHOR_NOT_FOUND})
			return
		}
		if count <= 0 {
			var authorCreateSuccess models.AutherCreateSuccess
			successMsg := utils.AUTHOR_NOT_FOUND
			authorCreateSuccess.Success = utils.BoolAddr(false)
			authorCreateSuccess.Message = &successMsg
			c.JSON(http.StatusCreated, authorCreateSuccess)
		} else {

			bookModl.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			bookModl.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			bookModl.ID = primitive.NewObjectID()
			tempuserId := bookModl.ID.Hex()
			bookModl.BookId = &tempuserId

			_, insertError := bookCollection.InsertOne(ctx, bookModl)
			if insertError != nil {
				msg := utils.ERROR_IN_BOOK_CREATION
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			}
			defer cancel()

			var authorCreateSuccess models.AutherCreateSuccess
			successMsg := utils.BOOK_ADD_SUCCESS
			authorCreateSuccess.Id = &tempuserId
			authorCreateSuccess.Success = utils.BoolAddr(true)
			authorCreateSuccess.Message = &successMsg
			c.JSON(http.StatusCreated, authorCreateSuccess)

		}
		//======================================

	}
}

func GetBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		id := c.Param("id")

		var foundBook models.Book
		var responseData models.GetBook

		bookCollection.FindOne(ctx, bson.M{"bookid": id}).Decode(&foundBook)
		defer cancel()

		if foundBook.BookId != nil {
			var foundAuthor models.Auther

			authorCollection.FindOne(ctx, bson.M{"authorid": foundBook.AuthorId}).Decode(&foundAuthor)
			var responseBook models.BookResponse
			responseBook.BookId = foundBook.AuthorId
			responseBook.CreatedAt = foundBook.CreatedAt
			responseBook.ID = foundBook.ID
			responseBook.ISBN = foundBook.ISBN
			responseBook.Name = foundBook.Name
			responseBook.Auther = &foundAuthor

			responseData.Success = utils.BoolAddr(true)
			responseData.Book = &responseBook
			c.JSON(http.StatusOK, responseData)
		} else {
			errorMsg := utils.BOOK_NOT_FOUND
			responseData.Success = utils.BoolAddr(false)
			responseData.Message = &errorMsg
			responseData.Book = nil
			c.JSON(http.StatusOK, responseData)

		}

	}
}

func GetAllBooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		page, err := strconv.Atoi(c.Param("page"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.INVALID_RQUEST})
			defer cancel()
			return
		}

		skip := (page - 1) * utils.MAX_PAGE_LIMIT
		findOptions := options.Find()
		findOptions.SetLimit(int64(utils.MAX_PAGE_LIMIT))
		findOptions.SetSkip(int64(skip))

		cursor, err := bookCollection.Find(ctx, bson.M{}, findOptions)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.INVALID_RQUEST})
			defer cancel()
			return
		}
		defer cancel()
		defer cursor.Close(ctx)
		var books []models.BookResponse
		var allList models.GetAllBooks

		for cursor.Next(ctx) {
			var book models.Book
			if err := cursor.Decode(&book); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": utils.INVALID_RQUEST})
				defer cancel()
				return
			}

			var foundAuthor models.Auther
			fmt.Println("BOOK ID : %s", *book.AuthorId)

			authorCollection.FindOne(ctx, bson.M{"userid": *book.AuthorId}).Decode(&foundAuthor)

			// fmt.Println("BOOK ID : %s", *book.AuthorId)

			var responseBook models.BookResponse
			responseBook.BookId = book.BookId
			responseBook.CreatedAt = book.CreatedAt
			responseBook.ID = book.ID
			responseBook.ISBN = book.ISBN
			responseBook.Name = book.Name
			responseBook.Auther = &foundAuthor

			books = append(books, responseBook)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.INVALID_RQUEST})
			defer cancel()
			return
		}
		allList.Author = &books
		allList.CurrentPage = &page
		nextPage := page + 1
		allList.NextPage = &nextPage
		c.JSON(http.StatusOK, allList)
	}
}
