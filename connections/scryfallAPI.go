/*
File		: scryfallAPI.go
Description	: File that deals with all the comunication with the Scryfall API.
*/

package connections

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"CardaliaAPI/models"
)

/*
Function	: Get card by name Scryfall
Description	: Given a card name with no spaces and commas, the function uses the ScryFall api to return a
card with all its information.

Parameters 	: cardName
Return     	: Card, error
*/
func GetCardByNameScryfall(cardName string) (models.Card, error) {
	resp, err := http.Get("https://api.scryfall.com/cards/named?exact=" + models.CleanCardName(cardName))
	var newCard models.Card
	if err != nil {
		return newCard, err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)

	// Convert response body to Card struct
	json.Unmarshal(bodyBytes, &newCard)
	return newCard, nil
}

/*
Function	: Get card by ID Scryfall
Description	: Given a cardID, the function uses the ScryFall api to return a card with all its information.
Parameters 	: cardID
Return     	: Card, error
*/
func GetCardByIDScryfall(ID string) (models.Card, error) {
	resp, err := http.Get("https://api.scryfall.com/cards/" + ID)
	var newCard models.Card
	if err != nil {
		return newCard, err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)

	// Convert response body to Card struct
	json.Unmarshal(bodyBytes, &newCard)
	return newCard, nil
}

/*
Function	: Get card by Uncompleted Scryfall
Description	: Given a partial card name, the function uses the ScryFall api to return a card with all its information.
This API consult can't generate error, only empty lists.

Parameters 	: uncompleted cardName
Return     	: cardName list
*/
func GetCardUncompletedScryfall(un_cardname string) []string {
	resp, err := http.Get("https://api.scryfall.com/cards/autocomplete?q=" + un_cardname)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer resp.Body.Close()

	var cardlist models.StringCardList
	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		json.Unmarshal(bodyBytes, &cardlist)
	}
	return cardlist.Cards
}

/*
Function	: Get card version Scryfall
Description	: Given a cardName, the function uses the ScryFall api to return all the paper versions of a card.
Parameters 	: cardName
Return     	: CardVersion list, error
*/
func GetCardVersionsScryfall(cardname string) ([]models.CardVersion, error) {
	var cardVersionsList []models.CardVersion
	resp, err := http.Get("https://api.scryfall.com/cards/search?order=released&q=%21%22" + cardname + "%22+include%3Aextras&unique=prints")
	if err != nil {
		return cardVersionsList, err
	}
	defer resp.Body.Close()

	var cardVersionsListRESP models.CardVersionsList
	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		json.Unmarshal(bodyBytes, &cardVersionsListRESP)
	}
	cardVersionsList = models.RemoveDigitalVersions(cardVersionsListRESP.Cards)
	return cardVersionsList, nil
}
