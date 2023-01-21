/*
File		: card.go
Description	: Model file to represent all the card-like objects and their related functions.
It also has implicid functions for the CardOwnership object.
*/

package models

import "strings"

// Object asociated to the CardOwnership table from the DB
type CardOwnership struct {
	CardID    uint   `gorm:"primary_key;auto_increment;not_null;" json:"card_id"`
	UserID    uint   `gorm:"not_null;foreignKey;" json:"user_id"`
	VersionID string `gorm:"not_null;" json:"version_id"`
	OracleID  string `gorm:"not_null;" json:"oracle_id"`
	Count     uint   `gorm:"not_null;" json:"count"`
	Extras    string `json:"extras"`
	Condi     string `json:"condi"`
}

// Used to get the info of a card from the Scryfall API and also used to send a card to the frontend.
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

// Used to get the info of the version of a card from the Scryfall API and also used to send a cardVersion to the frontend.
type CardVersion struct {
	Id              string    `json:"id"`
	Games           []string  `json:"games"`
	Set             string    `json:"set"`
	SetName         string    `json:"set_name"`
	CollectorNumber string    `json:"collector_number"`
	ImageURL        Image_url `json:"image_uris"`
}

// Object to get the card images
type Image_url struct {
	Small string `json:"small"`
	Large string `json:"large"`
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
Function	: Clear card name
Description	: Clear a card name by removing spaces and other punctuation hazards.
Self		: CardOwnership
Parameters 	:
Return     	: CardOwnership

*/
func (card *CardOwnership) SaveCard() (*CardOwnership, error) {
	var err error
	err = DB.Create(&card).Error
	if err != nil {
		return &CardOwnership{}, err
	}
	return card, nil
}

/*
Function	: Clear card name
Description	: Clear a card name by removing spaces and other punctuation hazards.
Parameters 	: cardName
Return     	: clean cardName
*/
func CleanCardName(exact_cardname string) string {
	cleaned_exact_cardname := strings.ReplaceAll(exact_cardname, " ", "")
	cleaned_exact_cardname = strings.ReplaceAll(cleaned_exact_cardname, ",", "")
	cleaned_exact_cardname = strings.ReplaceAll(cleaned_exact_cardname, "'", "")
	cleaned_exact_cardname = strings.ReplaceAll(cleaned_exact_cardname, "\n", "")
	cleaned_exact_cardname = strings.ReplaceAll(cleaned_exact_cardname, "\r", "")
	return cleaned_exact_cardname
}
