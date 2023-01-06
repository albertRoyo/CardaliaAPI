package models

import "strings"

type Card struct {
	Name            string    `json:"name"`
	Count           int       `json:"count"`
	ImageURL        Image_url `json:"image_uris"`
	ID              string    `json:"id"`
	VersionID       string    `json:"version_id"`
	OracleID        string    `json:"oracle_id"`
	Set             string    `json:"set"`
	SetName         string    `json:"set_name"`
	CollectorNumber string    `json:"collector_number"`
	Extras          string    `json:"extras"`
	Condi           string    `json:"condi"`
}
type Image_url struct {
	Small string `json:"small"`
	Large string `json:"large"`
}

type CardVersion struct {
	Id              string   `json:"id"`
	Games           []string `json:"games"`
	Set             string   `json:"set"`
	SetName         string   `json:"set_name"`
	CollectorNumber string   `json:"collector_number"`
}

type CardOwnership struct {
	CardID    uint   `gorm:"primary_key;auto_increment;not_null;" json:"card_id"`
	UserID    uint   `gorm:"not_null;foreignKey;" json:"user_id"`
	VersionID string `gorm:"not_null;" json:"version_id"`
	OracleID  string `gorm:"not_null;" json:"oracle_id"`
	Count     uint   `gorm:"not_null;" json:"count"`
	Extras    string `json:"extras"`
	Condi     string `json:"condi"`
}
type CardOwnerships struct {
	Collection []CardOwnership `json:"collection"`
}

/*
Removes spaces and other punctuation hazards.
TODO: See if there are some other puntuation problems:
*/
func CleanCardName(exact_cardname string) string {
	cleaned_exact_cardname := strings.ReplaceAll(exact_cardname, " ", "")
	cleaned_exact_cardname = strings.ReplaceAll(cleaned_exact_cardname, ",", "")
	cleaned_exact_cardname = strings.ReplaceAll(cleaned_exact_cardname, "'", "")
	cleaned_exact_cardname = strings.ReplaceAll(cleaned_exact_cardname, "\n", "")
	cleaned_exact_cardname = strings.ReplaceAll(cleaned_exact_cardname, "\r", "")
	return cleaned_exact_cardname
}

func (card *CardOwnership) SaveCard() (*CardOwnership, error) {
	var err error
	err = DB.Create(&card).Error
	if err != nil {
		return &CardOwnership{}, err
	}
	return card, nil
}

func GetVersionByCardID(card_id uint) string {
	var card CardOwnership
	card.CardID = card_id
	DB.First(&card)
	return card.VersionID
}

func GetCardIDByCardVersion(user_id_owner uint, versionID string) uint {
	var card CardOwnership
	card.VersionID = versionID
	card.UserID = user_id_owner
	DB.First(&card)
	return card.CardID
}
