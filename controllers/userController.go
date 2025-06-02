package controllers

import (
	"mas-diq/go-graphql/dto"
	"mas-diq/go-graphql/models"
	"mas-diq/go-graphql/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	res := schemas.Response{}

	var input dto.CreateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Name: input.Name, Email: input.Email}
	if err := models.CreateUserData(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res.Code = http.StatusOK
	res.Info = "User created successfully"
	res.Data = dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	c.JSON(http.StatusOK, res)
}

func GetUser(c *gin.Context) {
	res := schemas.Response{}
	id := c.MustGet("id").(uint64)

	var user models.User
	if err := models.GetOneUser(&user, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res.Code = http.StatusOK
	res.Info = "User retrieved successfully"
	res.Data = dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	c.JSON(http.StatusOK, res)
}

func GetListUser(c *gin.Context) {
	res := schemas.Response{}

	var user []models.User
	if err := models.GetListUser(&user, models.User{}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res.Code = http.StatusOK
	res.Info = "User retrieved successfully"
	res.Data = gin.H{
		"users": user,
	}
	c.JSON(http.StatusOK, res)
}

func UpdateUser(c *gin.Context) {
	res := schemas.Response{}
	id := c.MustGet("id").(uint64)

	var input dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := models.GetOneUser(&user, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Name = input.Name
	user.Email = input.Email
	if err := models.UpdateUserData(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res.Code = http.StatusOK
	res.Info = "User updated successfully"
	res.Data = dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	c.JSON(http.StatusOK, res)
}

func DeleteUser(c *gin.Context) {
	res := schemas.Response{}
	id := c.MustGet("id").(uint64)

	var user models.User
	if err := models.GetOneUser(&user, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.DeleteUser(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res.Code = http.StatusOK
	res.Info = "User deleted successfully"
	res.Data = dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	c.JSON(http.StatusOK, res)
}
