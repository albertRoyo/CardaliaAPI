/*
File		: token.go
Description	: File that deals with all the token related utilities.
*/

package token

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

/*
Function	: Generate Token
Description	: Generates a user token using JWT
Parameters 	: UserID
Return     	: Token, error
*/
func GenerateToken(user_id uint) (string, error) {
	// Decide the token life duration
	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		return "", err
	}
	// Define the 3 parameters of the token
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	return token.SignedString([]byte(os.Getenv("API_SECRET")))

}

/*
Function	: Token Validation
Description	: Validates the authority of a user token
Parameters 	: gin context -> request auth {token}
Return     	: error
*/
func TokenValid(c *gin.Context) error {
	// Get ths user token
	tokenString := extractToken(c)
	// Parse the token
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}
	return nil
}

/*
Function	: Extract Token ID
Description	: Extract userID from token
Parameters 	: gin context -> request auth {token}
Return     	: UserID, error
*/
func ExtractTokenID(c *gin.Context) (uint, error) {
	// Get ths user token
	tokenString := extractToken(c)
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}
	// Get the userID
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(uid), nil
	}
	return 0, nil
}

/*
Function	: Extract Token
Description	: Token extraction from gin context
Parameters 	: gin context -> request auth {token}
Return     	: Token
Private
*/
func extractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
