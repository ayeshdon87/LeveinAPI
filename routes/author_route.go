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

	incomingRoutes.GET(utils.V1_API_BASE_URL+"authors",
		controller.GetAllAuthers())

	incomingRoutes.PUT(utils.V1_API_BASE_URL+"author/update",
		controller.UpdateAuther())
}

func BookRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.POST(utils.V1_API_BASE_URL+"book/add",
		controller.AddBook())

	incomingRoutes.GET(utils.V1_API_BASE_URL+"book/:id",
		controller.GetBook())

	incomingRoutes.GET(utils.V1_API_BASE_URL+"books/:page",
		controller.GetAllBooks())

	incomingRoutes.PUT(utils.V1_API_BASE_URL+"book/update",
		controller.UpdateBook())
}
