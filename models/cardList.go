package models

// List with cards represented as strings
type StringCardList struct {
	Cards []string `json:"data"`
}

type CardList struct {
	Collection []Card `json:"collection"`
}

type CardVersionsList struct {
	Cards []CardVersion `json:"data"`
}

type Version struct {
	Set     string `json:"set"`
	SetName string `json:"set_name"`
	Number  string `json:"collector_number"`
}

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

func GetVersionNames(cardVersionsList []CardVersion) []Version {
	var versionList = []Version{}
	for _, cardVersion := range cardVersionsList {
		version := Version{Set: cardVersion.Set, SetName: cardVersion.SetName, Number: cardVersion.CollectorNumber}
		versionList = append(versionList, version)
	}
	return versionList
}

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
