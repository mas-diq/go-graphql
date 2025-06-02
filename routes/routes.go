package routes

import (
	"context"
	"mas-diq/go-graphql/config"
	"mas-diq/go-graphql/controllers"
	"mas-diq/go-graphql/graphql"
	"mas-diq/go-graphql/loaders"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
)

// Define a custom type for context keys
type contextKey string

const userLoaderKey contextKey = "userLoader"

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

	// In routes/routes.go
	r.POST("/graphql", func(c *gin.Context) {
		// Create new loader for each request
		loader := loaders.NewUserLoader(config.DB)
		ctx := context.WithValue(c.Request.Context(), userLoaderKey, loader)
		c.Request = c.Request.WithContext(ctx)
		h.ServeHTTP(c.Writer, c.Request)
	})

	return r
}
