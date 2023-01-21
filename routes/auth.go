/*
File		: auth.go
Description	: File that deals with all the HTTP requests that require authentification.
*/

package routes

import (
	"CardaliaAPI/models"
	"CardaliaAPI/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Function	: Change password
Description	: Changes user's password
Parameters 	: gin context -> request auth {token}

	-> request param {oldPassword, newPassword}

Return     	: message
*/
func ChangeUserPassword(c *gin.Context) {
	// Get ths userID that sends the request
	userID, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Bind the recived trade from gin.context
	var input models.UserChangePassword
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}
	u.User_id = userID

	// Change the user password
	err = u.ChangePassword(input.OldPassword, input.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

/*
Function	: Save collection (POST /user/collection)
Description	: Save the collection of the user.
Parameters 	: gin context -> request auth {token}
Return     	: message
*/
func SaveCollection(c *gin.Context) {
	// Get ths userID that sends the request
	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete the previous collection of the user from the DB
	if err := DeletePreviousCollection(user_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ownershipList := models.CardOwnershipList{}

	// Bind the recived collection from gin.context
	if err := c.ShouldBindJSON(&ownershipList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the recived collection to the DB
	if err := SaveUserCollection(ownershipList, user_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collection saved"})
}

/*
Function	: Get collection (GET /user/collection)
Description	: Get the collection of the user.
Parameters 	: gin context -> request auth {token}
Return     	: Card list
*/
func GetCollection(c *gin.Context) {
	// Get ths userID that sends the request
	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the user's collection from the DB
	collection, err := GetCollectionByUserID(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"collection": collection})
}

/*
Function	: Get all users collections by cardID(GET /users/collections/:cardname)
Description	: Get all the users collections that have a specific card.
Parameters 	: gin context	:card_id
Return     	: Collection list
*/
func GetAllUserCollectionsByCardId(c *gin.Context) {
	userIDAvoid, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userCollections, err := GetAllUserCollectionsByCardIdDB(userIDAvoid, c.Params.ByName("card_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_collections": userCollections})

}

/*
Function	: New Trade (POST /user/trade)
Description	: Start a new trade
Parameters 	: gin context -> request auth {token}

	-> request param {username, whatHeTrade, whatYouTrade, heChecked, youChecked}

Return     	: message
*/
func NewTrade(c *gin.Context) {
	// Get ths userID that sends the request
	user_id_origin, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Bind the recived trade from gin.context
	var holeTrade models.HoleTrade
	if err = c.ShouldBindJSON(&holeTrade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new trade
	if err = NewTradeDB(user_id_origin, holeTrade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "New trade offer made"})
}

/*
Function	: Modify Trade (PUT /user/trade)
Description	: Modify a parameter o a trade (username, whatHeTrade, whatYouTrade, heChecked, youChecked)
Parameters 	: gin context -> request auth {token}

	-> request param {username, whatHeTrade, whatYouTrade, heChecked, youChecked}

Return     	: message
*/
func ModifyTrade(c *gin.Context) {
	// Get ths userID that sends the request
	user_id_origin, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Bind the recived trade from gin.context
	var holeTrade models.HoleTrade
	if err = c.ShouldBindJSON(&holeTrade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Modify the trade
	if err = ModifyTradeDB(user_id_origin, holeTrade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trade updated"})
}

/*
Function	: Delete Trade (DELETE /user/trade)
Description	: Delete a trade between two users
Parameters 	: gin context -> request auth {token}
Return     	: message
*/
func DeleteTrade(c *gin.Context) {
	// Get ths userID that sends the request
	user_id1, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the target userID
	user_id2, err := models.GetUserIDByUsername(c.Params.ByName("username"))

	// Delete all trades that are not finished between users
	err = DeleteAllTradesBetweenUsersDB(user_id1, user_id2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trade deleted"})
}

/*
Function	: Get Trades (GET /user/trades)
Description	: Get all trades from the user.
Parameters 	: gin context -> request auth {token}
Return     	: Trade list
*/
func GetTrades(c *gin.Context) {
	// Get ths userID that sends the request
	user_id_origin, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get all trades that user participates in, including the finished ones.
	trades, err := GetTradesDB(user_id_origin)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trades": trades})
}
