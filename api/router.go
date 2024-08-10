package api

import (
	"user_medic/api/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "user_medic/api/docs"
)

// @title User Service API
// @version 1.0
// @description This is a sample server for a user service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:1001
// @BasePath /
func NewRouter(h *handler.Handler) *gin.Engine{
	router:=gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	user:=router.Group("/user")
	{
		user.POST("/register",h.Register)
		user.POST("/login",h.LoginUser)
		user.POST("/refresh-token",h.RefreshToken)
		user.GET("/profile/:email",h.GetUserProfile)
		user.PUT("/profile/update/:id/:email/:password/:first_name/:last_name/:date_of_birthday/:gender/:role",h.UpdateUserProfile)
		user.PUT("/logout",h.LogoutUser)
	}
	return router
}