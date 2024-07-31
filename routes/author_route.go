package routes

import (
	"github.com/ayeshdon87/LeveinAPI/controller"
	"github.com/ayeshdon87/LeveinAPI/utils"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.POST(utils.V1_API_BASE_URL+"author/add",
		controller.AddAuther())

	// incomingRoutes.GET(utils.V1_API_BASE_URL+"author/add",
	// func(c *gin.Context) {
	// 	var autherModl models.Auther

	// 	var signUpSuccess models.AutherCreateSuccess
	// 	tempuserId := autherModl.ID.Hex()
	// 	successMsg := utils.AUTHER_ADD_SUCCESS
	// 	signUpSuccess.Id = &tempuserId
	// 	signUpSuccess.Success = utils.BoolAddr(true)
	// 	signUpSuccess.Message = &successMsg
	// 	c.JSON(http.StatusCreated, signUpSuccess)
	// })
}
