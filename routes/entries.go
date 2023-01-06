package routes

import (
	"net/http"

	"github.com/albertRoyo/Cardalia.com/models"
	"github.com/gin-gonic/gin"
)

/*
Given an exact card name, this function calls GetCard to build a card struct with all its data
-IMPUT: 	gin context GET /card/:cardname
-RETURN:	Card(struct)
*/
func GetCardByName(c *gin.Context) {
	exact_card_name := models.CleanCardName(c.Params.ByName("cardname"))
	newCard, err := getCardByName(exact_card_name)

	if err != nil {
		c.IndentedJSON(http.StatusOK, newCard)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Card not found. Try again."})
	}
}

/*
Given an uncompleted card name, the function calls GetCardUncompleted and builds a list of card structs with the first 8 cards matches.
-IMPUT: 	gin context GET /:autocomplete
-RETURN:	Card(struct) list of 7 elements
*/
func GetCardsByName(c *gin.Context) {
	var cards = []string{}
	cards = getCardUncompleted(c.Params.ByName("autocomplete"))
	var showncards = []models.Card{}
	for index, card := range cards {
		var newCard models.Card
		newCard, err := getCardByName(card)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		showncards = append(showncards, newCard)
		//Change this value to get more or less Cards
		if index == 7 {
			break
		}
	}
	c.IndentedJSON(http.StatusOK, showncards)
}

/*
Get all the paper versions of a card.
-IMPUT: 	gin context GET /cards/versions/:cardname
-RETURN:	All the paper versions of a card
*/
func GetCardVersions(c *gin.Context) {
	cardVersionsList := getCardVersions(c.Params.ByName("cardname"))
	c.IndentedJSON(http.StatusOK, cardVersionsList)
}

/*
Get all the paper versions of a card.
-IMPUT: 	gin context GET /cards/versions/:cardname
-RETURN:	All the paper versions of a card
*/
func GetVersionNames(c *gin.Context) {
	cardVersionsList := getCardVersions(c.Params.ByName("cardname"))
	versionsNameList := models.GetVersionNames(cardVersionsList)
	c.IndentedJSON(http.StatusOK, versionsNameList)
}

/*
Get a version of a card.
-IMPUT: 	gin context GET /card/:set/:number
-RETURN:	A version of a card
*/
func GetCardVersion(c *gin.Context) {
	card := getCardVersion(c.Params.ByName("set"), c.Params.ByName("number"))

	c.IndentedJSON(http.StatusOK, card)
}

func GetUserCollectionByName(c *gin.Context) {
	collection, err := GetUserCollectionByNameDB(c.Params.ByName("username"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"user_collection": collection})
}
