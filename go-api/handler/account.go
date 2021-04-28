package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/whuangz/go-example/go-api/domain"
	"github.com/whuangz/go-example/go-api/middleware"
)

type accountHandler struct {
	service      domain.AccountService
	tokenService domain.TokenService
}

func NewAccountHandler(router *gin.Engine, service domain.AccountService, tokenService domain.TokenService) {
	h := &accountHandler{service: service,
		tokenService: tokenService}

	accountGroup := router.Group("/api/account")
	if gin.Mode() != gin.TestMode {
		//accountGroup.Use(middleware.Timeout(time.Duration(config.HANDLER_TIMEOUT), domain.NewServiceUnavailable()))
		accountGroup.GET("/me", middleware.AuthUser(tokenService), h.Me)
		accountGroup.POST("/signup", h.Signup)
		accountGroup.POST("/signin", h.Signin)
		accountGroup.POST("/signout", h.Signout)
		accountGroup.POST("/tokens", h.Tokens)
		accountGroup.POST("/image", h.Image)
		accountGroup.DELETE("/image", h.DeleteImage)
		accountGroup.PUT("/details", h.Details)
	} else {
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

	ctx := c.Request.Context()
	a, err := h.service.Get(ctx, uid)

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

	ctx := c.Request.Context()
	err := h.service.Signup(ctx, a)

	if err != nil {
		c.JSON(domain.Status(err), gin.H{
			"error": err,
		})
		return
	}

	token, err := h.tokenService.NewPairFromUser(ctx, a, "")
	if err != nil {
		c.JSON(domain.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": token,
	})
}

// Signin handler
func (h *accountHandler) Signin(c *gin.Context) {
	var req signInReq
	if ok := bindData(c, &req); !ok {
		return
	}

	a := &domain.Account{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := c.Request.Context()
	err := h.service.Signin(ctx, a)

	if err != nil {
		c.JSON(domain.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := h.tokenService.NewPairFromUser(ctx, a, "")
	if err != nil {
		c.JSON(domain.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tokens,
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
