package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/whuangz/go-example/go-api/domain"
)

type authHeader struct {
	AccessToken string `header:"Authorization"`
}

func AuthUser(s domain.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}

		if err := c.ShouldBindHeader(&h); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				var invalidArgs []domain.InvalidArgument

				for _, err := range errs {
					invalidArgs = append(invalidArgs, domain.InvalidArgument{
						err.Field(),
						err.Value().(string),
						err.Tag(),
						err.Param(),
					})
				}

				err := domain.NewBadRequest("Invalid request parameters. See invalidArgs")

				c.JSON(err.Status(), gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})
				c.Abort()
				return

			}

			err := domain.NewInternal()
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		accTokenHeader := strings.Split(h.AccessToken, "Bearer ")

		if len(accTokenHeader) < 2 {
			err := domain.NewAuthorization("Authorization header should with format `Bearer {token}`")

			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		acc, err := s.ValidateAccessToken(accTokenHeader[1])
		if err != nil {
			err := domain.NewAuthorization("Provided token is invalid")
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		c.Set("account", acc)

		c.Next()
	}
}
