package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/albertRoyo/CardaliaAPI/models"
)

/*
Given a card name with no spaces and commas, the function uses the ScryFall api to
create a new struct of the card.
-IMPUT: 	The exact cardname string
-RETURN:	Card(struct) |&&| True if the card was found, false otherwise
*/
func getCardByName(cardName string) (models.Card, error) {
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

func getCardByID(ID string) (models.Card, error) {
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
Given a partial card name, the function uses the ScryFall api to search for the first 20 cards matches.
-IMPUT: 	Partial card name string
-RETURN:	String card list
*/
func getCardUncompleted(un_cardname string) []string {
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
Get all the paper versions of a card.
-IMPUT: 	gin context GET /cards/versions/:cardname
-RETURN:	All the paper versions of a card
*/
func getCardVersions(cardname string) []models.CardVersion {
	resp, err := http.Get("https://api.scryfall.com/cards/search?order=released&q=%21%22" + cardname + "%22+include%3Aextras&unique=prints")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer resp.Body.Close()

	var cardVersionsListRESP models.CardVersionsList
	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		json.Unmarshal(bodyBytes, &cardVersionsListRESP)
	}
	cardVersionsList := models.RemoveDigitalVersions(cardVersionsListRESP.Cards)
	return cardVersionsList
}

/*
Get a version of a card.
-IMPUT: 	gin context GET /card/:set/:number
-RETURN:	A version of a card
*/
func getCardVersion(set string, number string) models.Card {
	resp, err := http.Get("https://api.scryfall.com/cards/" + set + "/" + number)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer resp.Body.Close()

	var newCard models.Card

	bodyBytes, _ := io.ReadAll(resp.Body)

	// Convert response body to Card struct
	json.Unmarshal(bodyBytes, &newCard)
	return newCard
}
