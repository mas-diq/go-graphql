package routes

import (
	"mas-diq/go-graphql/config"
	"mas-diq/go-graphql/controllers"
	"mas-diq/go-graphql/graphql"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// REST routes
	r.POST("/users", controllers.CreateUser)
	r.GET("/users/:id", controllers.GetUser)
	r.PUT("/users/:id", controllers.UpdateUser)
	r.DELETE("/users/:id", controllers.DeleteUser)

	// GraphQL route
	schema, _ := graphql.NewSchema(config.DB)
	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	r.POST("/graphql", func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})

	return r
}
