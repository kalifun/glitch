package gerr

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		// Content-Language
		contentLan := ctx.GetHeader("Content-Language")
		var lan string
		switch contentLan {
		case "en":
			lan = "en"
		case "zh":
			lan = "cn"
		default:
			lan = "en"
		}

		for _, err := range ctx.Errors {
			ge := &GError{}
			if errors.As(err, &ge) {
				ctx.JSON(http.StatusOK, ge.ToGinH(lan))
			} else {
				ctx.JSON(http.StatusBadGateway, gin.H{
					"error": map[string]string{
						"code":    "Unknow",
						"message": err.Error(),
					},
				})
			}
		}
	}
}
