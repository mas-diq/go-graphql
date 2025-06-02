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
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
}
