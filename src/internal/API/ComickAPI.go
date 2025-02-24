package API

import (
	"fmt"
	"strings"
)

const (
	ComickURL = "https://api.comick.io/v1.0/search/?page=1&limit=15&showall=false&q=Solo-Leveling&t=false"
)

type ComickTitleResponse struct {
	ID           uint64  `json:"id"`
	HID          string  `json:"hid"`  // Nome usado no URL.
	Slug         string  `json:"slug"` // Nome.
	BRating      float32 `json:"bayesian_rating"`
	Rating_count uint32  `json:"rating_count"`
	Description  string  `json:"desc"`
	Status       uint8   `json:"status"` // 2 -> Complete, 1 -> On Going,
	LastChapter  string  `json:"last_chapter"`
}

func SearchComickTitles(title string) (*ComickTitleResponse, error) {
	var res ComickTitleResponse // if res.Total > res.Limit {}
	search := ComickURL + "v1.0/search/?page=1&limit=15&showall=false&q=" + strings.ToLower(strings.ReplaceAll("-", " ", title)) + "&t=false"
	fmt.Println("Pesquisando: ", search)
	err := requestAndJsonfyResponse(search, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
