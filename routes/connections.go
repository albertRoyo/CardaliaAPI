/*
File		: connections.go
Description	: File that deals with all the communication with the data base.
*/

package routes

import (
	"errors"

	"CardaliaAPI/models"
	"CardaliaAPI/utils/token"

	"golang.org/x/crypto/bcrypt"
)

/*
Function	: Login Check
Description	: Checks in the DB if the users actually exists and if so, retruns a token.
Parameters 	: username, password
Return     	: email, token, error
*/
func LoginCheck(Username string, password string) (string, string, error) {
	var err error
	u := models.User{}
	// Get the user
	err = models.DB.Model(models.User{}).Where("Username = ?", Username).Take(&u).Error
	if err != nil {
		return "", "", err
	}
	// Check if the password is correct
	err = models.VerifyPassword(password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", "", err
	}
	// Generate token
	token, err := token.GenerateToken(u.User_id)
	if err != nil {
		return "", "", err
	}
	return u.Email, token, nil
}

/*
Function	: Get user collection by username
Description	: Get the user's collection from DB by his username
Parameters 	: username
Return     	: Collection, error
*/
func GetUserCollectionByNameDB(username string) ([]models.Card, error) {
	collection := []models.Card{}
	var user models.User
	// Get the user by its username
	err := models.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return collection, err
	}
	// Get the collection from the user
	collection, err = GetCollectionByUserID(user.User_id)
	return collection, nil
}

/*
Function	: Get all users collections by CardID
Description	: Get all the users collections from DB that have a specific card.
Parameters 	: cardID
Return     	: Collection list, error
*/
func GetAllUserCollectionsByCardIdDB(userIDAvoid uint, oracleID string) ([]models.UserCollection, error) {
	userCollections := []models.UserCollection{}
	// Get all the users that have the card
	users, err := getUsersWithCardDB(userIDAvoid, oracleID)
	if err != nil {
		return userCollections, err
	}
	// For each of those users
	for _, user := range users {
		// Get the user's collection from the DB
		collection, err := GetCollectionByUserID(user)
		if err != nil {
			return nil, err
		}

		var userCollection = models.UserCollection{}
		userCollection.Collection = collection

		// Get the collections username
		userCollection.Username, err = models.GetUsernameByUserID(user)
		if err != nil {
			return nil, err
		}
		// Append the collection to the collection list
		userCollections = append(userCollections, userCollection)
	}
	return userCollections, nil
}

/*
Function	: Delete the previous collection
Description	: Deletes from the DB all the cards in a user's collection.
Parameters 	: userID
Return     	: error
*/
func DeletePreviousCollection(user_id uint) error {
	err := models.DB.Where("user_id LIKE ?", user_id).Delete(&models.CardOwnership{}).Error
	if err != nil {
		return errors.New("Problem ocurred when deleting previous collection")
	}
	return nil
}

/*
Function	: Save users collection
Description	: Saves in the DB the user's collection.
Parameters 	: CardOwnership list, userID
Return     	: error
*/
func SaveUserCollection(ownershipList models.CardOwnershipList, userID uint) error {
	for _, card := range ownershipList.CardOwnerships {
		card.UserID = userID
		_, err := card.SaveCard()

		if err != nil {
			return err
		}
	}
	return nil
}

/*
Function	: Get user collection
Description	: Get from the DB the user's collection.
Parameters 	: userID
Return     	: Card list, error
*/
func GetCollectionByUserID(user_id uint) ([]models.Card, error) {

	var cardsByUserID []models.CardOwnership
	collection := []models.Card{}

	if err := models.DB.Where("user_id = ?", user_id).Find(&cardsByUserID).Error; err != nil {
		return collection, errors.New("User not found!")
	}

	for _, cardDB := range cardsByUserID {
		card, err := buildCard(cardDB)
		if err != nil {
			return collection, err
		}
		card.VersionID = card.ID
		collection = append(collection, card)
	}
	return collection, nil
}

/*
Function	: New Trade DB
Description	: Creates a new trade and store it to DB.
Parameters 	: userID, Trade list
Return     	: error
*/
func NewTradeDB(user_id_origin uint, holeTrade models.HoleTrade) error {
	// Get the userID of the owner of the card
	user_id_owner, err := models.GetUserIDByUsername(holeTrade.Username)
	if err != nil {
		return err
	}
	// For every cardOwnership the user has chosen
	for _, cardSelect := range holeTrade.WhatHeTrade {
		// Mount the Trade
		trade := models.Trade{}
		trade.User_id_origin = user_id_origin
		trade.User_id_owner = user_id_owner
		trade.VersionID = cardSelect.Card.VersionID
		trade.Extras = cardSelect.Card.Extras
		trade.Condi = cardSelect.Card.Condi
		trade.Card_select = cardSelect.Select
		trade.Status, err = getStatus(holeTrade.YouChecked, holeTrade.HeChecked, user_id_origin, user_id_owner)
		if err != nil {
			return err
		}
		// Save the new trade to the DB
		_, err := trade.CreateTrade()
		if err != nil {
			return err
		}
	}
	return nil
}

/*
Function	: Modify Trade
Description	: Modify a Trade in the DB.
Parameters 	: userID, Trade list
Return     	: error
*/
func ModifyTradeDB(user_id_origin uint, holeTrade models.HoleTrade) error {
	// Get the userID of the owner of the card
	user_id_owner, err := models.GetUserIDByUsername(holeTrade.Username)
	if err != nil {
		return err
	}
	// Delete all previous tredes between users that are not finished
	err = DeleteAllTradesBetweenUsersDB(user_id_origin, user_id_owner)
	if err != nil {
		return err
	}
	// For every cardOwnership the user has chosen
	for _, cardSelect := range holeTrade.WhatHeTrade {
		// Mount the Trade
		trade := models.Trade{}
		trade.User_id_origin = user_id_origin
		trade.User_id_owner = user_id_owner
		trade.VersionID = cardSelect.Card.VersionID
		trade.Extras = cardSelect.Card.Extras
		trade.Condi = cardSelect.Card.Condi
		trade.Card_select = cardSelect.Select
		// Get the status of the card
		trade.Status, err = getStatus(holeTrade.YouChecked, holeTrade.HeChecked, user_id_origin, user_id_owner)
		if err != nil {
			return err
		}
		println("Status he: ", trade.Status)
		// If the status of the trade is finished, delete the selection
		if trade.Status == 0 {
			err := deleteSelect(cardSelect, user_id_owner)
			if err != nil {
				return err
			}
		}

		_, err := trade.SaveTrade()
		if err != nil {
			return err
		}
	}
	// For every cardOwnership the other user has chosen
	for _, cardSelect := range holeTrade.WhatYouTrade {
		trade := models.Trade{}
		trade.User_id_origin = user_id_owner
		trade.User_id_owner = user_id_origin
		trade.VersionID = cardSelect.Card.VersionID
		trade.Extras = cardSelect.Card.Extras
		trade.Condi = cardSelect.Card.Condi
		trade.Card_select = cardSelect.Select
		trade.Status, err = getStatus(holeTrade.HeChecked, holeTrade.YouChecked, user_id_owner, user_id_origin)
		if err != nil {
			return err
		}
		println("Status you: ", trade.Status)
		if trade.Status == 0 {
			err := deleteSelect(cardSelect, user_id_origin)
			if err != nil {
				return err
			}
		}
		_, err := trade.SaveTrade()
		if err != nil {
			return err
		}
	}
	return nil
}

/*
Function	: Delete all trades between users
Description	: Delete all trades from the DB between two users.
Parameters 	: userID, userID
Return     	: error
*/
func DeleteAllTradesBetweenUsersDB(user1 uint, user2 uint) error {
	err := models.DB.Where("((user_id_origin = ? AND user_id_owner = ?) OR (user_id_origin = ? AND user_id_owner = ?)) AND status != ?", user1, user2, user2, user1, 0).Delete(&models.Trade{}).Error
	if err != nil {
		return err
	}
	return nil
}

/*
Function	: Get Trades
Description	: Get all trades from the user. This function builds a list of trades where in each element there are all the exchanges
with a specific username. The trades that are finished are put in another element with the same username

Parameters 	: userID
Return     	: Trade list, error
*/
func GetTradesDB(userAsking uint) ([]models.HoleTrade, error) {
	var trades []models.Trade
	tradeList := []models.HoleTrade{} //return
	tradeMap := make(map[string]models.HoleTrade)
	tradeMapFinished := make(map[string]models.HoleTrade)
	emptyTrades := make([]models.CardSelect, 0)
	// Get all the trades in where the user contributes
	if err := models.DB.Where("user_id_origin = ? OR user_id_owner = ?", userAsking, userAsking).Find(&trades).Error; err != nil {
		return tradeList, err
	}
	// For each of those trades
	for _, trade := range trades {
		var username string
		email := ""
		var emailProvi string
		// Decide the other user username
		if trade.User_id_origin == userAsking {
			username, _ = models.GetUsernameByUserID(trade.User_id_owner)
			emailProvi, _ = models.GetEmailByUserID(trade.User_id_owner)
		} else {
			username, _ = models.GetUsernameByUserID(trade.User_id_origin)
			emailProvi, _ = models.GetEmailByUserID(trade.User_id_origin)
		}
		var card models.Card
		var err error
		// If the trade is finished, the cards have been removed from users collections so we get it directly from Scryfall
		// Also, if the tarde is finished, we pass the email of the other user
		if trade.Status == 0 {
			card, err = GetCardByIDScryfall(trade.VersionID)
			card.Extras = trade.Extras
			card.Condi = trade.Condi
			email = emailProvi
			if err != nil {
				return tradeList, err
			}
		} else {
			// Get the card of the trade from the users collection
			card, err = getCardByParams(trade.User_id_owner, trade.VersionID, trade.Extras, trade.Condi)
			if err != nil {
				return tradeList, err
			}
		}

		// Mount the chechBoxes based on the satus
		var youChecked, heChecked = true, true
		if trade.Status != 0 {
			youChecked, heChecked = getCheks(trade, userAsking)
		}
		// Decide if the trade is finished to add it on a different map
		var holeTrade, ok = models.HoleTrade{}, false
		if trade.Status == 0 {
			holeTrade, ok = tradeMapFinished[username]
		} else {
			holeTrade, ok = tradeMap[username]
		}
		// Create a new map element
		if !ok {
			holeTrade = models.HoleTrade{
				Username:     username,
				Email:        email,
				WhatHeTrade:  emptyTrades,
				WhatYouTrade: emptyTrades,
				YouChecked:   youChecked,
				HeChecked:    heChecked,
			}
		}
		// Append to existing element
		if trade.User_id_origin == userAsking {
			holeTrade.WhatHeTrade = append(holeTrade.WhatHeTrade, models.CardSelect{
				Card:   card,
				Select: trade.Card_select,
			})

		} else {
			holeTrade.WhatYouTrade = append(holeTrade.WhatYouTrade, models.CardSelect{
				Card:   card,
				Select: trade.Card_select,
			})
		}
		// Modify existing element
		if trade.Status == 0 {
			tradeMapFinished[username] = holeTrade
		} else {
			tradeMap[username] = holeTrade
		}

	}

	// Concatenate the two maps into a list
	for _, v := range tradeMap {
		tradeList = append(tradeList, v)
	}
	for _, f := range tradeMapFinished {
		tradeList = append(tradeList, f)
	}

	return tradeList, nil
}

/*
Function	: Build Card
Description	: Build a card based on the info of cardOwnership
Parameters 	: CardOwnership
Return     	: Card, error
Private
*/
func buildCard(cardDB models.CardOwnership) (models.Card, error) {
	card, err := GetCardByIDScryfall(cardDB.VersionID)
	if err != nil {
		return card, err
	}

	card.Count = int(cardDB.Count)
	card.Extras = cardDB.Extras
	card.Condi = cardDB.Condi

	return card, nil
}

/*
Function	: Get Users with CardID
Description	: Get all users form the DB that have a specifict cardID
Parameters 	: CardID
Return     	: UserID list, error
Private
*/
func getUsersWithCardDB(userIDAvoid uint, oracleID string) ([]uint, error) {

	var cardOwnerships []models.CardOwnership
	err := models.DB.Table("card_ownerships").Where("oracle_id = ?", oracleID).Pluck("user_id", &cardOwnerships).Error
	if err != nil {
		return nil, err
	}

	var userIDs []uint
	for _, card := range cardOwnerships {
		if card.UserID != userIDAvoid {
			userIDs = append(userIDs, card.UserID)
		}
	}
	userIDs = models.RemoveDuplicate(userIDs)
	return userIDs, nil
}

/*
Function	: Get Status
Description	: Get the satatus of a trade based on the Checked values
Parameters 	: Check, Check, userID, userID
Return     	: Status(int), error
Private
*/
func getStatus(YouChecked bool, HeChecked bool, user_id_origin uint, user_id_owner uint) (int, error) {
	if YouChecked && HeChecked {
		return 0, nil
	} else if !YouChecked && HeChecked {
		return int(user_id_owner), nil
	} else if YouChecked && !HeChecked {
		return int(user_id_origin), nil
	} else {
		return -1, nil
	}
}

/*
Function	: Get Checks
Description	: Get the checks of a trade based on the trade Status
Parameters 	: Trade, UserID
Return     	: Check, Check
Private
*/
func getCheks(trade models.Trade, userAsking uint) (bool, bool) {
	if trade.Status == -1 {
		return false, false
	} else {
		if trade.Status == int(userAsking) {
			return true, false
		} else {
			return false, true
		}
	}
}

/*
Function	: Delete Select
Description	: Delete the selected cards from a user collection when a trade is finished
Parameters 	: Card select, userID
Return     	: error
Private
*/
func deleteSelect(cardSelect models.CardSelect, user_id_owner uint) error {
	cardOwnership := models.CardOwnership{}
	err := models.DB.Where("user_id = ? AND version_id = ? AND extras = ? AND condi = ?", user_id_owner, cardSelect.Card.VersionID, cardSelect.Card.Extras, cardSelect.Card.Condi).First(&cardOwnership).Error
	if err != nil {
		return err
	}
	cardOwnership.Count -= cardSelect.Select
	if cardOwnership.Count == 0 {
		// Delete the card from the CardOwnership table
		models.DB.Delete(&cardOwnership)
	} else {
		models.DB.Save(&cardOwnership)
	}
	return nil
}

/*
Function	: Get Card by parameters
Description	: Get a Card from the DB with a combinations of parameters that make it unique
Parameters 	: UserID, CardID, CardExtras, CardCondition
Return     	: Card, error
Private
*/
func getCardByParams(user_id uint, version_id string, extras string, condi string) (models.Card, error) {
	cardOwnership := models.CardOwnership{}
	card := models.Card{}
	err := models.DB.Where("user_id = ? AND version_id = ? AND extras = ? AND condi = ?", user_id, version_id, extras, condi).First(&cardOwnership).Error
	if err != nil {
		return card, err
	}

	card, err = GetCardByIDScryfall(cardOwnership.VersionID)
	card.VersionID = cardOwnership.VersionID
	card.Count = int(cardOwnership.Count)
	card.Extras = cardOwnership.Extras
	card.Condi = cardOwnership.Condi

	return card, nil

}
