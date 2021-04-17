package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/whuangz/go-example/go-api/domain"
)

type accountHandler struct {
	service domain.AccountService
}

func NewAccountHandler(router *gin.Engine, service domain.AccountService) {
	h := &accountHandler{service: service}

	accountGroup := router.Group("/api/account")
	{
		accountGroup.GET("/me", h.Me)
		accountGroup.POST("/signup", h.Signup)
		accountGroup.POST("/signin", h.Signin)
		accountGroup.POST("/signout", h.Signout)
		accountGroup.POST("/tokens", h.Tokens)
		accountGroup.POST("/image", h.Image)
		accountGroup.DELETE("/image", h.DeleteImage)
		accountGroup.PUT("/details", h.Details)
	}
}

// Me handler calls services for getting
// a user's details
func (h *accountHandler) Me(c *gin.Context) {
	account, exists := c.Get("account")

	if !exists {
		err := domain.NewAuthorization("unauthroized")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	uid := account.(*domain.Account).UID

	a, err := h.service.Get(c, uid)

	if err != nil {
		e := domain.NewNotFound("account", uid.String())

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": a,
	})

}

// Signup handler
func (h *accountHandler) Signup(c *gin.Context) {
	var req signupReq
	if ok := bindData(c, &req); !ok {
		return
	}

	a := &domain.Account{
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.service.Signup(c, a)

	if err != nil {
		c.JSON(domain.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "success",
	})
}

// Signin handler
func (h *accountHandler) Signin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's signin",
	})
}

// Signout handler
func (h *accountHandler) Signout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's signout",
	})
}

// Tokens handler
func (h *accountHandler) Tokens(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's tokens",
	})
}

// Image handler
func (h *accountHandler) Image(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's image",
	})
}

// DeleteImage handler
func (h *accountHandler) DeleteImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's deleteImage",
	})
}

// Details handler
func (h *accountHandler) Details(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's details",
	})
}
