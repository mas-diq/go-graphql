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

	// REST routes for users
	users := r.Group("users")
	{
		users.GET("", controllers.GetListUser)
		users.POST("", controllers.CreateUser)
		users.GET("/:id", controllers.GetUser)
		users.PUT("/:id", controllers.UpdateUser)
		users.DELETE("/:id", controllers.DeleteUser)
	}

	// REST routes for posts
	posts := r.Group("posts")
	{
		posts.POST("", controllers.CreatePost)
		posts.PUT("/:id", controllers.UpdatePost)
		posts.DELETE("/:id", controllers.DeletePost)
	}

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
