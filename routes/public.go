/*
File		: public.go
Description	: File that deals with all the HTTP requests that doesn't require authentification.
*/

package routes

import (
	"CardaliaAPI/connections"
	"net/http"

	"CardaliaAPI/models"

	"github.com/gin-gonic/gin"
)

/*
Function	: Register (POST /register)
Description	: Register a new user.
Parameters 	: gin context -> request params {username, password}
Return     	: message
*/
func Register(c *gin.Context) {
	//Get the parameters from the http request
	var input models.UserRegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}
	u.Username = input.Username
	u.Email = input.Email
	u.Password = input.Password

	// Encrypt password
	err := u.BeforeSave()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the new user in the DB
	_, err = u.SaveUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successfull"})
}

/*
Function	: Login (POST /login)
Description	: Login an existing user.
Parameters 	: gin context -> request params {username, password}
Return     	: token
*/
func Login(c *gin.Context) {
	//Get the parameters from the http request
	var input models.UserLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}
	u.Username = input.Username
	u.Password = input.Password

	// Generate a token
	email, token, err := connections.LoginCheck(u.Username, u.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"email": email, "token": token})
}

/*
Function	: Get cards by uncompleted cardname (GET /cards/:autocomplete)
Description	: Given an uncompleted card name, the function calls GetCardUncompleted and builds a list of card
structs with the first 8 cards matches.

Parameters 	: gin context 	:autocomplete
Return     	: Card list
*/
func GetCardsByName(c *gin.Context) {
	var cards = []string{}

	// Get the list of cards (strings) that match the search
	cards = connections.GetCardUncompletedScryfall(c.Params.ByName("autocomplete"))
	var showncards = []models.Card{}

	// For each card, get the card in Scryfall
	for index, card := range cards {
		var newCard models.Card
		newCard, err := connections.GetCardByNameScryfall(card)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		showncards = append(showncards, newCard)
		// Change this value to get more or less Cards
		if index == 7 {
			break
		}
	}
	c.IndentedJSON(http.StatusOK, showncards)
}

/*
Function	: Get card by cardname (GET /cards/versions/:cardname)
Description	: Get all the paper versions of a card.
Parameters 	: gin context	:cardname
Return     	: CardVersion list
*/
func GetCardVersions(c *gin.Context) {
	cardVersionsList, err := connections.GetCardVersionsScryfall(c.Params.ByName("cardname"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.IndentedJSON(http.StatusOK, cardVersionsList)
}

/*
Function	: Get a user collection by username (GET /user/collection/:username)
Description	: Get the collection of a user.
Parameters 	: gin context	:username
Return     	: Collection
*/
func GetUserCollectionByName(c *gin.Context) {
	collection, err := connections.GetUserCollectionByNameDB(c.Params.ByName("username"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"user_collection": collection})
}
