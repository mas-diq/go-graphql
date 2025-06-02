package graphql

import (
	"mas-diq/go-graphql/models"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":    &graphql.Field{Type: graphql.Int},
		"name":  &graphql.Field{Type: graphql.String},
		"email": &graphql.Field{Type: graphql.String},
		"createdAt": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				post := p.Source.(models.Post)
				return post.CreatedAt.Format("2006-01-02 15:04:05"), nil
			},
		},
		"updatedAt": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				post := p.Source.(models.Post)
				return post.UpdatedAt.Format("2006-01-02 15:04:05"), nil
			},
		},
	},
})

var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"id":       &graphql.Field{Type: graphql.Int},
		"title":    &graphql.Field{Type: graphql.String},
		"subtitle": &graphql.Field{Type: graphql.String},
		"image":    &graphql.Field{Type: graphql.String},
		"content":  &graphql.Field{Type: graphql.String},
		"status":   &graphql.Field{Type: graphql.String},
		"createdAt": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				post := p.Source.(models.Post)
				return post.CreatedAt.Format("2006-01-02 15:04:05"), nil
			},
		},
		"updatedAt": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				post := p.Source.(models.Post)
				return post.UpdatedAt.Format("2006-01-02 15:04:05"), nil
			},
		},
	},
})

func NewSchema(db *gorm.DB) (graphql.Schema, error) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, _ := p.Args["id"].(int)
					var user models.User
					db.First(&user, id)
					return user, nil
				},
			},

			// Post queries
			"post": &graphql.Field{
				Type: postType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					var post models.Post
					if err := db.First(&post, id).Error; err != nil {
						return nil, err
					}
					return post, nil
				},
			},
			"posts": &graphql.Field{
				Type: graphql.NewList(postType),
				Args: graphql.FieldConfigArgument{
					"status": &graphql.ArgumentConfig{Type: graphql.String},
					"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var posts []models.Post
					query := db.Model(&models.Post{})

					if status, ok := p.Args["status"].(string); ok {
						query = query.Where("status = ?", status)
					}

					if limit, ok := p.Args["limit"].(int); ok {
						query = query.Limit(limit)
					}

					if err := query.Find(&posts).Error; err != nil {
						return nil, err
					}
					return posts, nil
				},
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
}
