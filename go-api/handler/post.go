package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/whuangz/go-example/go-api/domain"
)

type postHandler struct {
	service domain.PostService
}

func NewPostHandler(router *gin.Engine, service domain.PostService) {
	handler := &postHandler{service}

	postGroup := router.Group("posts")
	{
		postGroup.GET("", handler.getPosts)
		//blogsGroup.POST("/", handler.createBlog)
	}
}

func (p *postHandler) getPosts(c *gin.Context) {
	posts, err := p.service.FindAll(c)

	if err != nil {
		c.JSON(404, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(200, posts)
}

// func (p *postHandler) createBlog(c *gin.Context) {

// 	postJSON := domain.Post{
// 		Title:   "test golang",
// 		Content: "test golang content",
// 	}

// 	// if err != nil {
// 	// 	c.JSON(400, gin.H{
// 	// 		"message": "Oops",
// 	// 	})
// 	// 	c.Abort()
// 	// 	return
// 	// }

// 	c.JSON(200, &postJSON)
// }
