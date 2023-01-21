/*
File		: cardList.go
Description	: Model file to represent all the list-like objects and their related functions.
*/

package models

// Represents CardOwnership list. Used to save the user collection.
type CardOwnershipList struct {
	CardOwnerships []CardOwnership `json:"collection"`
}

// Used to get the all the cards that match a uncompleated card search from ScryFall
type StringCardList struct {
	Cards []string `json:"data"`
}

// Used to get the all the versions of a single card from ScryFall
type CardVersionsList struct {
	Cards []CardVersion `json:"data"`
}

// Used to get all users collections.
type UserCollection struct {
	Username   string `json:"username"`
	Collection []Card `json:"collection"`
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
Function	: Remove digital versions
Description	: Removes all the digital version of a list of cardVersion to stay with the phisic versions.
Parameters 	: cardVersion list
Return     	: cardVersion list
*/
func RemoveDigitalVersions(cardVersionsList []CardVersion) []CardVersion {
	var cardVersionsListPaper = []CardVersion{}
	for index, card := range cardVersionsList {
		for _, game := range card.Games {
			if game == "paper" {
				cardVersionsListPaper = append(cardVersionsListPaper, cardVersionsList[index])
			}
		}
	}
	return cardVersionsListPaper
}

/*
Function	: Remove duplicate
Description	: Removes all duplicates from a UserID list.
Parameters 	: UserID list
Return     	: UserID list
*/
func RemoveDuplicate(intSlice []uint) []uint {
	allKeys := make(map[uint]bool)
	list := []uint{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
