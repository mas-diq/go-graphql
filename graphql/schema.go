package graphql

import (
	"mas-diq/go-graphql/loaders"
	"mas-diq/go-graphql/models"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func NewSchema(db *gorm.DB) (graphql.Schema, error) {
	// Define base types without relationships
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":    &graphql.Field{Type: graphql.Int},
			"name":  &graphql.Field{Type: graphql.String},
			"email": &graphql.Field{Type: graphql.String},
		},
	})

	postType := graphql.NewObject(graphql.ObjectConfig{
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
		},
	})

	// Add relationships using AddFieldConfig
	userType.AddFieldConfig("posts", &graphql.Field{
		Type: graphql.NewList(postType),
		Args: graphql.FieldConfigArgument{
			"limit": &graphql.ArgumentConfig{Type: graphql.Int},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			user := p.Source.(models.User)
			var posts []models.Post
			query := db.Where("created_by = ?", user.ID)

			if limit, ok := p.Args["limit"].(int); ok {
				query = query.Limit(limit)
			}

			if err := query.Find(&posts).Error; err != nil {
				return nil, err
			}
			return posts, nil
		},
	})

	postType.AddFieldConfig("author", &graphql.Field{
		Type: userType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// Get loader from context
			loader := p.Context.Value("userLoader").(*loaders.UserLoader)

			post := p.Source.(models.Post)
			users, err := loader.Load(p.Context, []uint{post.CreatedBy})
			if err != nil {
				return nil, err
			}
			return users[0], nil
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					var user models.User
					if err := db.First(&user, id).Error; err != nil {
						return nil, err
					}
					return user, nil
				},
			},
			"post": &graphql.Field{
				Type: postType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					var post models.Post
					if err := db.Preload("User").First(&post, id).Error; err != nil {
						return nil, err
					}
					return post, nil
				},
			},
			"posts": &graphql.Field{
				Type: graphql.NewList(postType),
				Args: graphql.FieldConfigArgument{
					"status":   &graphql.ArgumentConfig{Type: graphql.String},
					"limit":    &graphql.ArgumentConfig{Type: graphql.Int},
					"authorId": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var posts []models.Post
					query := db

					if status, ok := p.Args["status"].(string); ok {
						query = query.Where("status = ?", status)
					}

					if authorId, ok := p.Args["authorId"].(int); ok {
						query = query.Where("created_by = ?", authorId)
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
