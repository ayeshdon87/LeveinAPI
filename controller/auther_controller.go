package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/ayeshdon87/LeveinAPI/database"
	"github.com/ayeshdon87/LeveinAPI/models"
	"github.com/ayeshdon87/LeveinAPI/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()

var authorCollection *mongo.Collection = database.OpentCollection(database.Client,
	utils.AUTHER_TABLE)

func AddAuther() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var autherModl models.Auther

		err := c.BindJSON(&autherModl)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}
		//validate request format
		validationError := validate.Struct(autherModl)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.INVALID_RQUEST})
			defer cancel()
			return
		}

		autherModl.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		autherModl.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		autherModl.ID = primitive.NewObjectID()
		tempuserId := autherModl.ID.Hex()
		autherModl.UserId = &tempuserId

		_, insertError := authorCollection.InsertOne(ctx, autherModl)
		if insertError != nil {
			msg := utils.ERROR_IN_AUTHER_CREATION
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		defer cancel()

		var authorCreateSuccess models.AutherCreateSuccess
		successMsg := utils.AUTHER_ADD_SUCCESS
		authorCreateSuccess.Id = &tempuserId
		authorCreateSuccess.Success = utils.BoolAddr(true)
		authorCreateSuccess.Message = &successMsg
		c.JSON(http.StatusCreated, authorCreateSuccess)
	}
}

func GetAuther() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		id := c.Param("id")
		//log.Panicln(`ID--> %s `, id)

		var foundAuthor models.Auther
		var responseData models.GetAuthor
		//var responseUser models.IsUserAvailableResponse
		authorCollection.FindOne(ctx, bson.M{"userid": id}).Decode(&foundAuthor)
		defer cancel()

		if foundAuthor.UserId != nil {
			responseData.Success = utils.BoolAddr(true)
			responseData.Author = &foundAuthor
		} else {
			errorMsg := utils.AUTHOR_NOT_FOUND
			responseData.Success = utils.BoolAddr(false)
			responseData.Message = &errorMsg
			responseData.Author = nil

		}

		c.JSON(http.StatusCreated, responseData)

	}
}
