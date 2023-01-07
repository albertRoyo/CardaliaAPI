package routes

import (
	"errors"

	"github.com/albertRoyo/CardaliaAPI/models"
	"github.com/albertRoyo/CardaliaAPI/utils/token"
	"golang.org/x/crypto/bcrypt"
)

func LoginCheck(Username string, password string) (string, error) {
	var err error
	u := models.User{}
	err = models.DB.Model(models.User{}).Where("Username = ?", Username).Take(&u).Error

	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	token, err := token.GenerateToken(u.User_id)

	if err != nil {
		return "", err
	}
	return token, nil
}

func deletePreviousCollection(user_id uint) error {
	err := models.DB.Where("user_id LIKE ?", user_id).Delete(&models.CardOwnership{}).Error
	if err != nil {
		return errors.New("Problem ocurred when deleting previous collection")
	}
	return nil
}

func saveUserCollection(collection models.CardOwnerships, userID uint) error {
	for _, card := range collection.Collection {
		card.UserID = userID
		_, err := card.SaveCard()

		if err != nil {
			return err
		}
	}
	return nil
}

func GetCollectionByUserID(user_id uint) ([]models.Card, error) {

	var cardsByUserID []models.CardOwnership
	collection := []models.Card{}

	if err := models.DB.Where("user_id = ?", user_id).Find(&cardsByUserID).Error; err != nil {
		return collection, errors.New("User not found!")
	}

	for _, cardDB := range cardsByUserID {
		card, err := prepareCollection(cardDB)
		if err != nil {
			return collection, err
		}
		card.VersionID = card.ID
		collection = append(collection, card)
	}
	return collection, nil
}

func prepareCollection(cardDB models.CardOwnership) (models.Card, error) {
	card, err := getCardByID(cardDB.VersionID)
	if err != nil {
		return card, err
	}

	card.Count = int(cardDB.Count)
	card.Extras = cardDB.Extras
	card.Condi = cardDB.Condi

	return card, nil
}

type UserCollection struct {
	Username   string        `json:"username"`
	Collection []models.Card `json:"collection"`
}

func GetUserCollectionByNameDB(username string) ([]models.Card, error) {
	collection := []models.Card{}
	var user models.User
	err := models.DB.Where("username = ?", username).First(&user).Error

	if err != nil {
		return collection, err
	}
	collection, err = GetCollectionByUserID(user.User_id)
	return collection, nil

}

func GetAllUserCollectionsByCardIdDB(userIDAvoid uint, oracleID string) ([]UserCollection, error) {
	userCollections := []UserCollection{}
	users, err := getUsersWithCardDB(userIDAvoid, oracleID)
	if err != nil {
		return userCollections, err
	}
	for _, user := range users {

		// Get the user's collection from the DB
		collection, err := GetCollectionByUserID(user)
		if err != nil {
			return nil, err
		}

		var userCollection = UserCollection{}
		userCollection.Collection = collection

		userCollection.Username, err = models.GetUsernameByUserID(user)
		if err != nil {
			return nil, err
		}

		userCollections = append(userCollections, userCollection)
	}
	return userCollections, nil

}

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

func NewTradeDB(user_id_origin uint, holeTrade models.TradeJSON) error {
	user_id_owner, err := models.GetUserIDByUsername(holeTrade.Username)
	if err != nil {
		return err
	}
	for _, cardSelect := range holeTrade.WhatHeTrade {
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
		_, err := trade.CreateTrade()
		if err != nil {
			return err
		}
	}
	return nil
}

func ModifyTradeDB(user_id_origin uint, holeTrade models.TradeJSON) error {
	user_id_owner, err := models.GetUserIDByUsername(holeTrade.Username)
	if err != nil {
		return err
	}

	err = DeleteAllTradesBetweenUsersDB(user_id_origin, user_id_owner)
	if err != nil {
		return err
	}

	for _, cardSelect := range holeTrade.WhatHeTrade {
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
		if trade.Status == 0 {
			err := DeleteSelect(cardSelect, user_id_owner)
			if err != nil {
				return err
			}
		}

		_, err := trade.SaveTrade()
		if err != nil {
			return err
		}
	}

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
		if trade.Status == 0 {
			err := DeleteSelect(cardSelect, user_id_origin)
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

func DeleteSelect(cardSelect models.CardJSON, user_id_owner uint) error {
	cardOwnership := models.CardOwnership{}
	err := models.DB.Where("user_id = ? AND version_id = ? AND extras = ? AND condi = ?", user_id_owner, cardSelect.Card.VersionID, cardSelect.Card.Extras, cardSelect.Card.Condi).First(&cardOwnership).Error
	if err != nil {
		return err
	}

	cardOwnership.Count -= cardSelect.Select
	models.DB.Save(&cardOwnership)
	return nil
}

func DeleteAllTradesBetweenUsersDB(user1 uint, user2 uint) error {
	err := models.DB.Where("((user_id_origin = ? AND user_id_owner = ?) OR (user_id_origin = ? AND user_id_owner = ?)) AND status != ?", user1, user2, user2, user1, 0).Delete(&models.Trade{}).Error
	if err != nil {
		return err
	}
	return nil
}

func GetTradesDB(userAsking uint) ([]models.TradeJSON, error) {
	var trades []models.Trade
	tradeList := []models.TradeJSON{} //return
	tradeMap := make(map[string]models.TradeJSON)
	tradeMapFinished := make(map[string]models.TradeJSON)
	emptyTrades := make([]models.CardJSON, 0)

	if err := models.DB.Where("user_id_origin = ? OR user_id_owner = ?", userAsking, userAsking).Find(&trades).Error; err != nil {
		return tradeList, err
	}

	for _, trade := range trades {
		var username string
		if trade.User_id_origin == userAsking {
			username, _ = models.GetUsernameByUserID(trade.User_id_owner)
		} else {
			username, _ = models.GetUsernameByUserID(trade.User_id_origin)
		}
		card, err := GetCardByParams(trade.User_id_owner, trade.VersionID, trade.Extras, trade.Condi)
		if err != nil {
			return tradeList, err
		}
		var youChecked, heChecked = true, true
		if trade.Status != 0 {
			youChecked, heChecked = models.GetCheks(trade, userAsking)
		}

		var tradeJSON, ok = models.TradeJSON{}, false
		if trade.Status == 0 {
			tradeJSON, ok = tradeMapFinished[username]
		} else {
			tradeJSON, ok = tradeMap[username]
		}

		if !ok {
			tradeJSON = models.TradeJSON{
				Username:     username,
				WhatHeTrade:  emptyTrades,
				WhatYouTrade: emptyTrades,
				YouChecked:   youChecked,
				HeChecked:    heChecked,
			}
		}

		if trade.User_id_origin == userAsking {
			tradeJSON.WhatHeTrade = append(tradeJSON.WhatHeTrade, models.CardJSON{
				Card:   card,
				Select: trade.Card_select,
			})

		} else {
			tradeJSON.WhatYouTrade = append(tradeJSON.WhatYouTrade, models.CardJSON{
				Card:   card,
				Select: trade.Card_select,
			})
		}
		if trade.Status == 0 {
			tradeMapFinished[username] = tradeJSON
		} else {
			tradeMap[username] = tradeJSON
		}

	}

	for _, v := range tradeMap {
		tradeList = append(tradeList, v)
	}
	for _, f := range tradeMapFinished {
		tradeList = append(tradeList, f)
	}

	return tradeList, nil
}

/*

func GetTradesDB(userAsking uint) ([]models.TradeJSON, error) {
	var trades []models.Trade
	tradeList := []models.TradeJSON{} //return
	tradeMap := make(map[string]models.TradeJSON)

	// Get all t
	if err := models.DB.Where("user_id_origin = ? OR user_id_owner = ?", userAsking, userAsking).Find(&trades).Error; err != nil {
		return tradeList, err
	}

	for _, trade := range trades {
		var username string
		var err error
		if trade.User_id_origin == userAsking {
			username, err = models.GetUsernameByUserID(trade.User_id_owner)
		} else {
			username, err = models.GetUsernameByUserID(trade.User_id_origin)
		}
		if err != nil {
			return tradeList, err
		}
		card, err := GetCardByParams(trade.User_id_owner, trade.VersionID, trade.Extras, trade.Condi)
		if err != nil {
			return tradeList, err
		}
		tradeJSON, ok := tradeMap[username]
		usernameA := username + "$F"
		tradeJSONA, okA := tradeMap[usernameA]
		emptyTrades := []models.CardJSON{}
		youChecked, heChecked := getCheks(trade)

		if ok && (trade.User_id_origin == userAsking) && (trade.Status != 0) {
			// afegir a username, he
			tradeJSON.WhatHeTrade = append(tradeJSON.WhatHeTrade, models.CardJSON{
				Card:   card,
				Select: trade.Card_select,
			})
			tradeMap[username] = tradeJSON
		} else if !ok && (trade.User_id_origin == userAsking) && (trade.Status != 0) {
			// crear nou username, he
			tradeMap[username] = models.TradeJSON{
				Username: username,
				WhatHeTrade: []models.CardJSON{
					{
						Card:   card,
						Select: trade.Card_select,
					},
				},
				WhatYouTrade: emptyTrades,
				YouChecked:   youChecked,
				HeChecked:    heChecked,
			}
		} else if okA && (trade.User_id_origin == userAsking) && !(trade.Status != 0) {
			// afegir a usernameA, he
			tradeJSONA.WhatHeTrade = append(tradeJSONA.WhatHeTrade, models.CardJSON{
				Card:   card,
				Select: trade.Card_select,
			})
			tradeMap[usernameA] = tradeJSONA
		} else if !okA && (trade.User_id_origin == userAsking) && !(trade.Status != 0) {
			// crear nou usernameA, he
			tradeMap[usernameA] = models.TradeJSON{
				Username: usernameA,
				WhatHeTrade: []models.CardJSON{
					{
						Card:   card,
						Select: trade.Card_select,
					},
				},
				WhatYouTrade: emptyTrades,
				YouChecked:   youChecked,
				HeChecked:    heChecked,
			}
		} else if ok && !(trade.User_id_origin == userAsking) && (trade.Status != 0) {
			// afegir a username, you
			tradeJSON.WhatYouTrade = append(tradeJSON.WhatYouTrade, models.CardJSON{
				Card:   card,
				Select: trade.Card_select,
			})
			tradeMap[username] = tradeJSON
		} else if !ok && !(trade.User_id_origin == userAsking) && (trade.Status != 0) {
			// crear nou username, you
			tradeMap[username] = models.TradeJSON{
				Username: username,
				WhatYouTrade: []models.CardJSON{
					{
						Card:   card,
						Select: trade.Card_select,
					},
				},
				WhatHeTrade: emptyTrades,
				YouChecked:  youChecked,
				HeChecked:   heChecked,
			}
		} else if okA && !(trade.User_id_origin == userAsking) && !(trade.Status != 0) {
			// afegir a usernameA, you
			tradeJSONA.WhatYouTrade = append(tradeJSONA.WhatYouTrade, models.CardJSON{
				Card:   card,
				Select: trade.Card_select,
			})
			tradeMap[usernameA] = tradeJSONA
		} else if !okA && !(trade.User_id_origin == userAsking) && !(trade.Status != 0) {
			// crear nou usernameA, you
			tradeMap[usernameA] = models.TradeJSON{
				Username: usernameA,
				WhatYouTrade: []models.CardJSON{
					{
						Card:   card,
						Select: trade.Card_select,
					},
				},
				WhatHeTrade: emptyTrades,
				YouChecked:  youChecked,
				HeChecked:   heChecked,
			}
		}
	}

	for _, tradeJSON := range tradeMap {
		tradeList = append(tradeList, tradeJSON)
	}

	return tradeList, nil
}
*/

func GetCardByParams(user_id uint, version_id string, extras string, condi string) (models.Card, error) {
	cardOwnership := models.CardOwnership{}
	card := models.Card{}
	err := models.DB.Where("user_id = ? AND version_id = ? AND extras = ? AND condi = ?", user_id, version_id, extras, condi).First(&cardOwnership).Error
	if err != nil {
		return card, err
	}

	card, err = getCardByID(cardOwnership.VersionID)
	card.VersionID = cardOwnership.VersionID
	card.Count = int(cardOwnership.Count)
	card.Extras = cardOwnership.Extras
	card.Condi = cardOwnership.Condi

	return card, nil

}
