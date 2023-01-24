/*
File		: card.go
Description	: Model file to represent all the card-like objects and their related functions.
It also has functions for the CardOwnership object.
*/

package models

import (
	"strings"

	"gorm.io/gorm"
)

// Object asociated to the CardOwnership table from the DB
type CardOwnership struct {
	CardID    uint   `gorm:"primary_key;auto_increment;not_null;" json:"card_id"`
	User_id   uint   `gorm:"not_null;foreignKey;" json:"user_id"`
	VersionID string `gorm:"not_null;" json:"version_id"` ///ID to identyfy a card version
	OracleID  string `gorm:"not_null;" json:"oracle_id"`  //ID to identify a card. All versions of a single card have the same OracleID
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
Function	: Save Card
Description	: This function updates the card or creates a new one.

	The search parameters thet make a unique combination are: UserID, VersionID, Extras, Condi

Self		: CardOwnership
Parameters 	:
Return     	: CardOwnership
*/
func (card *CardOwnership) SaveCard() (*CardOwnership, error) {
	// Find existing card by unique combination of fields
	existingCard := &CardOwnership{}
	err := DB.Where("user_id = ? AND version_id = ? AND extras = ? AND condi = ?", card.User_id, card.VersionID, card.Extras, card.Condi).First(existingCard).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Update existing card or create new one
	if existingCard.CardID != 0 {
		existingCard.Count = card.Count
		existingCard.VersionID = card.VersionID
		existingCard.OracleID = card.OracleID
		existingCard.Extras = card.Extras
		existingCard.Condi = card.Condi
		DB.Save(existingCard)
		return existingCard, nil
	} else {
		DB.Create(card)
		return card, nil
	}
}

/*
Function	: Get cardOwnership by CardID
Description	: Get the CardOwnership from the DB with primary key CardID.
Parameters 	: CardID
Return     	: VersionID
*/
func GetCardOwnershipByCardID(cardId uint) (CardOwnership, error) {
	cardOwnership := CardOwnership{}
	err := DB.First(&cardOwnership, cardId).Error
	if err != nil {
		return cardOwnership, err
	}
	return cardOwnership, nil
}

/*
Function	: Get card ID by parameters
Description	: Get a CardID from the DB with a unique combinations of parameters (without primary key).
Parameters 	: UserID, CardID, CardExtras, CardCondition
Return     	: Card, error
Private
*/
func GetCardIDByParams(userId uint, versionId string, extras string, condi string) (uint, error) {
	cardOwnership := CardOwnership{}
	err := DB.Where("user_id = ? AND version_id = ? AND extras = ? AND condi = ?", userId, versionId, extras, condi).First(&cardOwnership).Error
	if err != nil {
		return 0, err
	}
	return cardOwnership.CardID, nil
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
