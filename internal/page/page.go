package page

import (
	"encoding/json"
	"log"
)

type Page struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Article      string `json:"article"`
	Posted       string `json:"posted"`
	Author       string `json:"author"`
	CommentCount int    `json:"comment_count"`
	Rating       int    `json:"rating"`
}

func (p *Page) GetJSON() []byte {
	js, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("%v", err.Error())
	}

	return js
}
