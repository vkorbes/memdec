package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Deck struct {
	DeckID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	LastShownIndex int           `json:"lastshownindex" bson:"lastshownindex"`
	Cards          []Card        `json:"cards" bson:"cards"`
}

type Card struct {
	Suit   string `json:"suit" bson:"suit"`
	Number string `json:"number" bson:"number"`
}
