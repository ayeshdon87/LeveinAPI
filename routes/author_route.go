package routes

import (
	"github.com/ayeshdon87/LeveinAPI/controller"
	"github.com/ayeshdon87/LeveinAPI/utils"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.POST(utils.V1_API_BASE_URL+"author/add",
		controller.AddAuther())

	incomingRoutes.GET(utils.V1_API_BASE_URL+"author/:id",
		controller.GetAuther())

	incomingRoutes.GET(utils.V1_API_BASE_URL+"authors/:page",
		controller.GetAllAuthers())
}
