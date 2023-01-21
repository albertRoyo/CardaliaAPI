/*
File		: authMiddlewares.go
Description	: File used to create a middleware for user authentification.
*/

package middlewares

import (
	"net/http"

	"CardaliaAPI/utils/token"

	"github.com/gin-gonic/gin"
)

/*
Function	: JWT Auth Middleware
Description	: Checks if the token sent by the user is valid
Parameters 	: username, password
Return     	: token, error
*/
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
