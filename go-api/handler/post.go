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

	postGroup := router.Group("/api/post")
	{
		postGroup.GET("", handler.getPosts)
		postGroup.POST("", handler.createPost)
		postGroup.GET("/:post_id", handler.getPostByID)
		postGroup.PATCH("/:post_id", handler.updatePost)
		postGroup.DELETE("/:post_id", handler.deletePost)
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

func (p *postHandler) createPost(c *gin.Context) {
	var req domain.Post
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(422, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	err = p.service.Save(c, &req)

	if err != nil {
		c.JSON(domain.Status(err), gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(200, &req)
}

func (p *postHandler) getPostByID(c *gin.Context) {
	if postId, ok := getPathInt(c, "post_id"); ok {

		post, err := p.service.FindByID(c, int32(postId))

		if err != nil {
			c.JSON(domain.Status(err), gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		c.JSON(200, post)
	}
}

func (p *postHandler) updatePost(c *gin.Context) {
	if postId, ok := getPathInt(c, "post_id"); ok {

		var req domain.Post
		if ok := bindData(c, &req); !ok {
			return
		}
		err := p.service.Update(c, int32(postId), &req)

		if err != nil {
			c.JSON(domain.Status(err), gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		c.JSON(200, &req)
	}
}

func (p *postHandler) deletePost(c *gin.Context) {
	if postId, ok := getPathInt(c, "post_id"); ok {
		err := p.service.Delete(c, int32(postId))

		if err != nil {
			c.JSON(domain.Status(err), gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		c.JSON(200, "success deleted")
	}
}
