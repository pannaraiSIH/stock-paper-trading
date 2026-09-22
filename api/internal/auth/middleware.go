package auth

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pannaraiSIH/stock-paper-trading/internal/response"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, ok := bearerToken(ctx)
		if !ok {
			response.Unauthorized(ctx, "unauthorized")
			ctx.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			response.Unauthorized(ctx, "unauthorized")
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			response.Unauthorized(ctx, "unauthorized")
			ctx.Abort()
			return
		}

		ctx.Set("userID", claims.UserID)
		ctx.Next()
	}
}

func bearerToken(ctx *gin.Context) (string, bool) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], true
		}
		return "", false
	}

	if strings.EqualFold(ctx.GetHeader("Upgrade"), "websocket") {
		if token := ctx.Query("token"); token != "" {
			return token, true
		}
	}

	return "", false
}
