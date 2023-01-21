/*
File		: trade.go
Description	: Model file to represent all the trade-like objects and their related functions.
It also has implicid functions for the Trade object.
*/

package models

// Trade DB object.
type Trade struct {
	Trade_id       uint   `json:"trade_id"`
	User_id_origin uint   `json:"user_id_origin"`
	User_id_owner  uint   `json:"user_id_owner"`
	VersionID      string `gorm:"not_null;" json:"version_id"`
	Extras         string `json:"extras"`
	Condi          string `json:"condi"`
	Card_select    uint   `json:"card_select"`
	Status         int    `json:"status"`
}

// Object that represents all the trades a user has.
type HoleTrade struct {
	Username     string       `json:"username"`
	Email        string       `json:"email"`
	WhatHeTrade  []CardSelect `json:"whatHeTrade"`
	WhatYouTrade []CardSelect `json:"whatYouTrade"`
	YouChecked   bool         `json:"youChecked"`
	HeChecked    bool         `json:"heChecked"`
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
Return     	: Trade
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
Return     	: Trade, error
*/
func (trade *Trade) SaveTrade() (*Trade, error) {
	if DB.Model(&trade).Where("user_id_origin = ? AND user_id_owner = ? AND version_id = ? AND extras = ? AND condi = ? AND status != ?",
		trade.User_id_origin, trade.User_id_owner, trade.VersionID, trade.Extras, trade.Condi, 0).Updates(&trade).RowsAffected == 0 {
		if err := DB.Create(&trade).Error; err != nil {
			return &Trade{}, err
		}
	}
	return trade, nil
}
