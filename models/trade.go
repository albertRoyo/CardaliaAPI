package models

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

type TradeJSON struct {
	Username     string     `json:"username"`
	WhatHeTrade  []CardJSON `json:"whatHeTrade"`
	WhatYouTrade []CardJSON `json:"whatYouTrade"`
	YouChecked   bool       `json:"youChecked"`
	HeChecked    bool       `json:"heChecked"`
}

type CardJSON struct {
	Card   Card `json:"card"`
	Select uint `json:"select"`
}

func (trade *Trade) CreateTrade() (*Trade, error) {
	if err := DB.Create(&trade).Error; err != nil {
		return &Trade{}, err
	}
	return trade, nil
}

func (trade *Trade) SaveTrade() (*Trade, error) {

	if DB.Model(&trade).Where("user_id_origin = ? AND user_id_owner = ? AND version_id = ? AND extras = ? AND condi = ? AND status != ?",
		trade.User_id_origin, trade.User_id_owner, trade.VersionID, trade.Extras, trade.Condi, 0).Updates(&trade).RowsAffected == 0 {
		if err := DB.Create(&trade).Error; err != nil {
			return &Trade{}, err
		}
	}
	return trade, nil
}

func GetCheks(trade Trade, userAsking uint) (bool, bool) {
	if trade.Status == -1 {
		return false, false
	} else {
		if trade.User_id_origin == userAsking && trade.Status == int(trade.User_id_origin) {
			return true, false
		} else {
			return false, true
		}
	}
}
