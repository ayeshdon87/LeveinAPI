package controller

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/ayeshdon87/LeveinAPI/database"
	"github.com/ayeshdon87/LeveinAPI/models"
	"github.com/ayeshdon87/LeveinAPI/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

		var foundAuthor models.Auther
		var responseData models.GetAuthor

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

		c.JSON(http.StatusOK, responseData)

	}
}

func GetAllAuthers() gin.HandlerFunc {
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

		cursor, err := authorCollection.Find(ctx, bson.M{}, findOptions)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.INVALID_RQUEST})
			defer cancel()
			return
		}
		defer cancel()
		defer cursor.Close(ctx)
		var authers []models.Auther
		var allList models.GetAllAuthor

		for cursor.Next(ctx) {
			var auther models.Auther
			if err := cursor.Decode(&auther); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": utils.INVALID_RQUEST})
				defer cancel()
				return
			}
			authers = append(authers, auther)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.INVALID_RQUEST})
			defer cancel()
			return
		}
		allList.Author = &authers
		allList.CurrentPage = &page
		nextPage := page + 1
		allList.NextPage = &nextPage
		c.JSON(http.StatusOK, allList)
	}
}

func UpdateAuther() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var autherModl models.AutherUpdate

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
		autherModl.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		var updateObj primitive.D
		updateObj = append(updateObj, bson.E{Key: "firstname", Value: *autherModl.FirstName})
		updateObj = append(updateObj, bson.E{Key: "lastname", Value: *autherModl.LastName})
		UpdatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updatedat", Value: UpdatedAt})
		upsert := true
		objID, err := primitive.ObjectIDFromHex(*autherModl.UserId)
		filter := bson.M{"_id": objID}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, updateError := authorCollection.UpdateOne(ctx, filter, bson.D{
			{Key: "$set", Value: updateObj},
		},
			&opt)

		if updateError != nil {
			msg := utils.ERROR_IN_AUTHER_CREATION
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		defer cancel()
		var authorCreateSuccess models.AutherCreateSuccess

		if result.MatchedCount == 0 {
			successMsg := utils.AUTHOR_NOT_FOUND
			authorCreateSuccess.Success = utils.BoolAddr(false)
			authorCreateSuccess.Message = &successMsg
		} else {
			successMsg := utils.AUTHER_UPDATE_SUCCESS
			authorCreateSuccess.Success = utils.BoolAddr(true)
			authorCreateSuccess.Message = &successMsg
		}

		c.JSON(http.StatusCreated, authorCreateSuccess)
	}
}
