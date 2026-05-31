package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/luponetn/noitrex/utils"
)

func AuthMiddleware(jwtAccessSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			slog.Error("access_token cookie missing", "error", err)
			utils.Unauthorized(c)
			c.Abort()
			return
		}

		claims, err := utils.VerifyJwt(tokenString, jwtAccessSecret)
		if err != nil {
			slog.Error("access_token validation failed", "error", err.Error())
			utils.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set("operatorId", claims.OperatorID)

		c.Next()
	}
}
