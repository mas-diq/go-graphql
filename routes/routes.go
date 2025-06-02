package routes

import (
	"mas-diq/go-graphql/config"
	"mas-diq/go-graphql/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	userController := &controllers.UserController{DB: config.DB}

	// REST routes
	r.POST("/users", userController.CreateUser)
	r.GET("/users/:id", userController.GetUser)
	r.PUT("/users/:id", userController.UpdateUser)
	r.DELETE("/users/:id", userController.DeleteUser)

	// GraphQL route
	// schema, _ := graphql.NewSchema(config.DB)
	// h := handler.New(&handler.Config{
	// 	Schema: &schema,
	// 	Pretty: true,
	// })

	// r.POST("/graphql", func(c *gin.Context) {
	// 	h.ServeHTTP(c.Writer, c.Request)
	// })

	return r
}
