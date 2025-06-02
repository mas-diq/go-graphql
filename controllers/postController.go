package controllers

import (
	"mas-diq/go-graphql/dto"
	"mas-diq/go-graphql/models"
	"mas-diq/go-graphql/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	res := schemas.Response{}

	var input dto.CreatePostRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := models.Post{
		Title:     input.Title,
		Subtitle:  input.Subtitle,
		Image:     input.Image,
		Content:   input.Content,
		Status:    models.PostStatus(input.Status),
		CreatedBy: input.CreatedBy,
	}

	if err := models.CreatePostData(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res.Code = http.StatusOK
	res.Info = "Post created successfully"
	res.Data = dto.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Subtitle:  post.Subtitle,
		Image:     post.Image,
		Content:   post.Content,
		Status:    string(post.Status),
		CreatedBy: post.CreatedBy,
	}
	c.JSON(http.StatusOK, res)
}

func UpdatePost(c *gin.Context) {
	res := schemas.Response{}
	id := c.MustGet("id").(uint64)

	var input dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post models.Post
	if err := models.GetOnePost(&post, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.Title = input.Title
	post.Subtitle = input.Subtitle
	post.Image = input.Image
	post.Content = input.Content
	post.Status = models.PostStatus(input.Status)
	if err := models.UpdatePostData(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res.Code = http.StatusOK
	res.Info = "Post updated successfully"
	res.Data = dto.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Subtitle:  post.Subtitle,
		Image:     post.Image,
		Content:   post.Content,
		Status:    string(post.Status),
		CreatedBy: post.CreatedBy,
	}
	c.JSON(http.StatusOK, res)
}

func DeletePost(c *gin.Context) {
	res := schemas.Response{}
	id := c.MustGet("id").(uint64)

	var post models.Post
	if err := models.GetOnePost(&post, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.DeletePost(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res.Code = http.StatusOK
	res.Info = "Post deleted successfully"
	res.Data = nil
	c.JSON(http.StatusOK, res)
}
