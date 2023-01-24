/*
File		: trade.go
Description	: Model file to represent all the trade-like objects and their related functions.
It also has implicid functions for the Trade object.
*/

package models

// Trade DB object.
type Trade struct {
	TradeID      uint `gorm:"primary_key;auto_increment;not_null;" json:"trade_id"`
	UserIdOrigin uint `gorm:"not_null;" json:"user_id_origin"`
	UserIdOwner  uint `gorm:"not_null;" json:"user_id_owner"`
	CardID       uint `gorm:"not_null;foreignKey;" json:"card_id"`
	//VersionID    string `gorm:"not_null;" json:"version_id"`
	//Extras       string `json:"extras"`
	//Condi        string `json:"condi"`
	CardSelect uint `json:"card_select"`
	Status     int  `json:"status"`
	// -1 if both users dont want to finish
	// 0 if both users want to finish
	// UserIdOrigin if the user asking the card/s wants to finish
	// UserIdowner if the user owning the card/s wants to finish
}

// Object that represents all the trades a user has with another user.
type HoleTrade struct {
	Username     string       `json:"username"`     // The other user username
	Email        string       `json:"email"`        // The other user email
	WhatHeTrade  []CardSelect `json:"whatHeTrade"`  // The cards that the other user gives
	WhatYouTrade []CardSelect `json:"whatYouTrade"` // The cards that the other user gives
	YouChecked   bool         `json:"youChecked"`   // True if the user wants to finish the trade
	HeChecked    bool         `json:"heChecked"`    // True if the other user wants to finish the trade
}

// Object that represents the number of selections of a traded card.
type CardSelect struct {
	Card   Card `json:"card"`
	Select uint `json:"select"`
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
Function	: Create Trade
Description	: Store a new trade to the DB.
Self		: Trade
Parameters 	:
Return     	: *Trade, error
*/
func (trade *Trade) CreateTrade() (*Trade, error) {
	if err := DB.Create(&trade).Error; err != nil {
		return &Trade{}, err
	}
	return trade, nil
}

/*
Function	: Save Trade
Description	: Modify a trade from the DB and if not found, create a new one.
Self		: Trade
Parameters 	:
Return     	: *Trade, error
*/
func (trade *Trade) SaveTrade() (*Trade, error) {
	// Find
	if DB.Model(&trade).Where("user_id_origin = ? AND user_id_owner = ? AND card_id = ? ",
		trade.UserIdOrigin, trade.UserIdOwner, trade.CardID).Updates(&trade).RowsAffected == 0 {
		if err := DB.Create(&trade).Error; err != nil {
			return &Trade{}, err
		}
	}
	return trade, nil
}
