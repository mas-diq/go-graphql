package graphql

import (
	"fmt"
	"mas-diq/go-graphql/loaders"
	"mas-diq/go-graphql/models"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

// NewSchema creates and returns a new GraphQL schema.
// It defines the types, relationships, and resolvers for querying data via GraphQL.
// The 'db' parameter is a GORM database instance used for data retrieval.
func NewSchema(db *gorm.DB) (graphql.Schema, error) {
	// --- Define base GraphQL object types without relationships first ---

	// userType defines the GraphQL 'User' object.
	// It specifies the fields available for a user.
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User", // Name of the type in the GraphQL schema
		Fields: graphql.Fields{
			"id":    &graphql.Field{Type: graphql.Int},    // User's unique identifier
			"name":  &graphql.Field{Type: graphql.String}, // User's name
			"email": &graphql.Field{Type: graphql.String}, // User's email address
		},
	})

	// postType defines the GraphQL 'Post' object.
	// It specifies the fields available for a post.
	postType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Post", // Name of the type in the GraphQL schema
		Fields: graphql.Fields{
			"id":       &graphql.Field{Type: graphql.Int},    // Post's unique identifier
			"title":    &graphql.Field{Type: graphql.String}, // Post's title
			"subtitle": &graphql.Field{Type: graphql.String}, // Post's subtitle
			"image":    &graphql.Field{Type: graphql.String}, // URL or path to the post's image
			"content":  &graphql.Field{Type: graphql.String}, // Main content of the post
			"status":   &graphql.Field{Type: graphql.String}, // Status of the post (e.g., "published", "draft")
			"createdAt": &graphql.Field{ // Post's creation timestamp
				Type: graphql.String, // Exposed as a formatted string
				// Resolve function for 'createdAt' field.
				// It formats the timestamp from the database model into a readable string.
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// p.Source contains the parent object, which is a 'Post' model instance here.
					if post, ok := p.Source.(models.Post); ok { // Type assertion
						return post.CreatedAt.Format("2006-01-02 15:04:05"), nil // Standard Go time formatting
					}
					if post, ok := p.Source.(*models.Post); ok { // Type assertion for pointer
						return post.CreatedAt.Format("2006-01-02 15:04:05"), nil
					}
					return nil, fmt.Errorf("could not cast source to Post or *Post for createdAt")
				},
			},
		},
	})

	// --- Add relationships between types using AddFieldConfig ---
	// This is done after base types are defined to avoid circular dependency issues during initialization.

	// Add 'posts' field to 'userType'.
	// This defines a one-to-many relationship: a User can have multiple Posts.
	userType.AddFieldConfig("posts", &graphql.Field{
		Type: graphql.NewList(postType), // The field returns a list of 'Post' types
		Args: graphql.FieldConfigArgument{
			// 'limit' argument allows clients to specify the maximum number of posts to retrieve.
			"limit": &graphql.ArgumentConfig{Type: graphql.Int},
		},
		// Resolve function for 'posts' field on User.
		// It fetches posts created by the specific user.
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// p.Source is the parent User object.
			user, ok := p.Source.(models.User) // Type assertion
			if !ok {
				userPtr, okPtr := p.Source.(*models.User)
				if !okPtr {
					return nil, fmt.Errorf("could not cast source to User or *User for user.posts resolver")
				}
				user = *userPtr
			}

			var posts []models.Post
			// Start building the GORM query to find posts where 'created_by' matches the user's ID.
			query := db.Where("created_by = ?", user.ID)

			// Apply 'limit' if provided in the GraphQL query arguments.
			if limit, ok := p.Args["limit"].(int); ok && limit > 0 {
				query = query.Limit(limit)
			}

			// Execute the query.
			if err := query.Find(&posts).Error; err != nil {
				return nil, err // Return error if database query fails
			}
			return posts, nil // Return the list of found posts
		},
	})

	// Add 'author' field to 'postType'.
	// This defines a many-to-one relationship: a Post has one Author (User).
	postType.AddFieldConfig("author", &graphql.Field{
		Type: userType, // The field returns a 'User' type
		// Resolve function for 'author' field on Post.
		// It fetches the user who created the post, utilizing a DataLoader for efficiency.
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// Retrieve the UserLoader instance from the GraphQL context.
			// DataLoaders are typically initialized per request and passed via context.
			loader, ok := p.Context.Value("userLoader").(*loaders.UserLoader)
			if !ok {
				// Fallback or error if loader is not found; for this example, we'll try a direct fetch.
				// In a real app, you might return an error:
				// return nil, fmt.Errorf("userLoader not found in context")

				// Fallback: Direct database query if loader is not available (less efficient for multiple author lookups)
				post, okPost := p.Source.(models.Post)
				if !okPost {
					postPtr, okPostPtr := p.Source.(*models.Post)
					if !okPostPtr {
						return nil, fmt.Errorf("could not cast source to Post or *Post for post.author resolver (fallback)")
					}
					post = *postPtr
				}
				var author models.User
				if err := db.First(&author, post.CreatedBy).Error; err != nil {
					return nil, err
				}
				return &author, nil // Return a pointer to the user, GORM often works well with pointers for single results
			}

			// p.Source is the parent Post object.
			post, okPost := p.Source.(models.Post)
			if !okPost {
				postPtr, okPostPtr := p.Source.(*models.Post)
				if !okPostPtr {
					return nil, fmt.Errorf("could not cast source to Post or *Post for post.author resolver (loader path)")
				}
				post = *postPtr
			}

			// Use the DataLoader to fetch the user.
			// DataLoader batches and caches requests, preventing N+1 query problems.
			// It expects a slice of keys (user IDs in this case) and returns a slice of users and a slice of errors.
			users, err := loader.Load(p.Context, []uint{post.CreatedBy})
			if err != nil { // Check for error from the loader
				return nil, err
			}
			if len(users) == 0 {
				return nil, fmt.Errorf("user not found by loader for ID: %d", post.CreatedBy)
			}
			return users[0], nil // Return the first user from the result (as we loaded by a single ID)
		},
	})

	// --- Define the Root Query type ---
	// queryType is the entry point for all GraphQL read operations.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query", // Standard name for the root query type
		Fields: graphql.Fields{
			// 'user' query field: Fetches a single user by their ID.
			"user": &graphql.Field{
				Type: userType, // Specifies that this query returns a 'User'
				Args: graphql.FieldConfigArgument{
					// 'id' argument is required to identify the user.
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				// Resolve function for the 'user' query.
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, _ := p.Args["id"].(int) // Get 'id' argument
					var user models.User
					// Fetch the user from the database using GORM.
					if err := db.First(&user, id).Error; err != nil {
						return nil, err // Return error if user not found or DB error
					}
					return &user, nil // Return the found user (pointer often preferred for GORM results)
				},
			},
			// 'post' query field: Fetches a single post by its ID.
			"post": &graphql.Field{
				Type: postType, // Specifies that this query returns a 'Post'
				Args: graphql.FieldConfigArgument{
					// 'id' argument is required to identify the post.
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				// Resolve function for the 'post' query.
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, _ := p.Args["id"].(int) // Get 'id' argument
					var post models.Post
					// Fetch the post from the database.
					// 'Preload("User")' tells GORM to also fetch the associated User (author)
					// to optimize and avoid a separate query if the 'author' field is requested.
					// Note: This preloads the 'User' field on the 'Post' GORM model,
					//       not necessarily what the GraphQL 'author' resolver for postType does by default.
					//       The GraphQL 'author' resolver will still run, potentially using a DataLoader.
					//       If the DataLoader is smart or if the user is already on the `post.User` struct,
					//       it can use this preloaded data.
					if err := db.Preload("User").First(&post, id).Error; err != nil {
						return nil, err // Return error if post not found or DB error
					}
					return &post, nil // Return the found post
				},
			},
			// 'posts' query field: Fetches a list of posts, with optional filters.
			"posts": &graphql.Field{
				Type: graphql.NewList(postType), // Specifies that this query returns a list of 'Post'
				Args: graphql.FieldConfigArgument{
					// 'status' argument to filter posts by their status.
					"status": &graphql.ArgumentConfig{Type: graphql.String},
					// 'limit' argument to restrict the number of posts returned.
					"limit": &graphql.ArgumentConfig{Type: graphql.Int},
					// 'authorId' argument to filter posts by the author's ID.
					"authorId": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				// Resolve function for the 'posts' query.
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var posts []models.Post
					query := db.Model(&models.Post{}) // Start with a base GORM query for Post model

					// Apply 'status' filter if provided.
					if status, ok := p.Args["status"].(string); ok && status != "" {
						query = query.Where("status = ?", status)
					}

					// Apply 'authorId' filter if provided.
					if authorId, ok := p.Args["authorId"].(int); ok && authorId > 0 {
						query = query.Where("created_by = ?", authorId)
					}

					// Apply 'limit' if provided.
					if limit, ok := p.Args["limit"].(int); ok && limit > 0 {
						query = query.Limit(limit)
					}

					// Execute the query to find all matching posts.
					if err := query.Find(&posts).Error; err != nil {
						return nil, err // Return error if database query fails
					}
					return posts, nil // Return the list of found posts
				},
			},
		},
	})

	// --- Create and return the GraphQL schema ---
	// The schema is configured with the root query type.
	// Mutations and Subscriptions would also be added here if defined.
	return graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType, // Set the root query type
		// Mutation: mutationType, // Example if you had mutations
	})
}
