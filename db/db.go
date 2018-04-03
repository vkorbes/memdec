package db

import (
	"github.com/ellenkorbes/memdec/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Init(arg string) *mgo.Session {
	session, err := mgo.Dial(arg)
	if err != nil {
		panic(err)
	}
	return session
}

func AddDeck(db *mgo.Session, deck models.Deck) error {
	return db.DB("memdec").C("decks").Insert(deck)
}

func GetAllDecks(db *mgo.Session) ([]models.Deck, error) {
	list := []models.Deck{}
	err := db.DB("memdec").C("decks").Find(bson.M{}).All(&list)
	return list, err
}

func GetDeck(db *mgo.Session, id bson.ObjectId) (models.Deck, error) {
	deck := models.Deck{}
	err := db.DB("memdec").C("decks").FindId(id).One(&deck)
	return deck, err
}

func IsUnique(db *mgo.Session, deck bson.ObjectId) (bool, error) {
	c := db.DB("memdec").C("decks")
	count, err := c.Find(bson.M{"DeckID": deck}).Limit(1).Count()
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func IncrementLastShown(db *mgo.Session, deck models.Deck) error {
	deckCheck := models.Deck{}
	plusOne := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"lastshownindex": 1}},
		ReturnNew: true,
	}
	_, err := db.DB("memdec").C("decks").FindId(deck.DeckID).Apply(plusOne, &deckCheck)
	if err != nil {
		return err
	}
	return nil
}
