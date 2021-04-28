package handler

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/whuangz/go-example/go-api/domain"
)

func getPathInt(c *gin.Context, name string) (int, bool) {
	val := c.Params.ByName(name)
	r, err := strconv.Atoi(val)
	if err != nil || val == "" {
		c.JSON(404, gin.H{
			"error": "Missing param id",
		})
		c.Abort()
		return 0, false
	}

	return r, true
}

func bindData(c *gin.Context, req interface{}) bool {
	if c.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", c.FullPath())

		err := domain.NewUnsupportedMediaType(msg)

		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		c.Abort()
		return false
	}
	// Bind incoming json to struct and check for validation errors
	if err := c.ShouldBind(req); err != nil {
		log.Printf("Error binding data: %+v\n", err)

		if errs, ok := err.(validator.ValidationErrors); ok {
			// could probably extract this, it is also in middleware_auth_user
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
			return false
		}

		// later we'll add code for validating max body size here!

		// if we aren't able to properly extract validation errors,
		// we'll fallback and return an internal server error
		fallBack := domain.NewInternal()

		c.JSON(fallBack.Status(), gin.H{"error": fallBack})
		return false
	}

	return true
}
