# Go REST & GraphQL API Server

This project provides a dual-interface API server with both REST and GraphQL endpoints for managing users and posts. Built with Gin and GORM, it features efficient data loading and relational queries.

## Project Structure
```bash
.
├── config
│   └── database.go       # Database configuration
├── controllers
│   ├── postController.go # Post REST handlers
│   └── userController.go # User REST handlers
├── dto
│   ├── postDto.go        # Post data transfer objects
│   └── userDto.go        # User data transfer objects
├── go.mod                # Go dependencies
├── go.sum                # Dependency checksums
├── graphql
│   └── schema.go         # GraphQL schema definition
├── loaders
│   └── loaders.go        # DataLoader implementation
├── main.go               # Entry point
├── models
│   ├── post.go           # Post model
│   └── user.go           # User model
├── routes
│   └── routes.go         # Route configuration
└── schemas
    └── schemas.go        # Database schema setup
```

## Features
- REST API with CRUD operations
- GraphQL endpoint with queries
- N+1 prevention using DataLoader
- Relational data handling (users ↔ posts)
- Request batching and caching

## Prerequisites
- Go 1.16+
- PostgreSQL
- gqlgen (go install github.com/graphql-go/graphql@latest)

## Setup

### Clone the repository:
```bash
git clone https://github.com/mas-diq/go-graphql.git
cd go-graphql
```

### Install dependencies:
```bash
go mod download
```

### Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your database credentials
```

### Run migrations:
```bash
go run schemas/schemas.go
```

### Start the server:
```bash
go run main.go
```

## REST API Endpoints
### User Routes
| Method | Endpoint   | Description     |
|--------|------------|-----------------|
| GET    | /users     | List all users  |
| POST   | /users     | Create new user |
| GET    | /users/:id | Get user by ID  |
| PUT    | /users/:id | Update user     |
| DELETE | /users/:id | Delete user     |

### Post Routes
| Method | Endpoint   | Description     |
|--------|------------|-----------------|
| POST   | /posts     | Create new post |
| PUT    | /posts/:id | Update post     |
| DELETE | /posts/:id | Delete post     |

## GraphQL API
### Endpoint
| Method | Endpoint   | Description     |
|--------|------------|-----------------|
| POST   | /graphql   | Graphql queries |
 

### Example Queries
```graphql
# Get user with posts
query GetUserWithPosts($userId: Int!) {
  user(id: $userId) {
    id
    name
    email
    posts(limit: 3) {
      id
      title
      createdAt
    }
  }
}

# Get posts with authors
query GetPosts {
  posts(limit: 5, status: "published") {
    id
    title
    createdAt
    author {
      id
      name
    }
  }
}
```

Query Variables
```json
{
  "userId": 1
}
```

## Data Loader Implementation
The GraphQL resolver uses DataLoader to batch user requests when resolving post authors:

```go
// In routes/routes.go
loader := loaders.NewUserLoader(config.DB)
ctx := context.WithValue(c.Request.Context(), userLoaderKey, loader)
```

## Configuration
Edit config/database.go for DB settings:

```go
func ConnectDB() *gorm.DB {
  // ... PostgreSQL connection setup
}
```

## Testing
Use curl to test REST endpoints:

```bash
# Create user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John", "email":"john@example.com"}'
```

## Get user
```bash
curl http://localhost:8080/users/1
```

## Development
Regenerate GraphQL schema (if modified):

```bash
go run github.com/graphql-go/graphql generate
```

## Format code:
```bash
gofmt -w .
```

## Deployment
Build binary:

```go
go build -o api-server main.go
./api-server
```

## Dependencies
- Gin (Web framework)
- GORM (ORM)
- gqlgen (GraphQL implementation)
- PostgreSQL driver