package ctrl

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"path"
	"time"

	"github.com/ellenkorbes/memdec/db"
	"github.com/ellenkorbes/memdec/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Controller struct {
	DB *mgo.Session
}

func NewController(db *mgo.Session) *Controller {
	return &Controller{
		DB: db,
	}
}

func (c *Controller) ListAllDecks(response http.ResponseWriter, request *http.Request) {
	items, err := db.GetAllDecks(c.DB)
	if err != nil {
		panic(nil)
	}
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(items)
}

func (c *Controller) Info(response http.ResponseWriter, request *http.Request) {
	id := path.Base(request.URL.Path)
	deck := c.fetchDeck(id)
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(&deck)
}

func (c *Controller) fetchDeck(id string) models.Deck {
	if !bson.IsObjectIdHex(id) {
		panic("invalid ID")
	}
	deck, err := db.GetDeck(c.DB, bson.ObjectIdHex(id))
	if err != nil {
		panic(err)
	}
	return deck
}

func (c *Controller) NextCard(response http.ResponseWriter, request *http.Request) {
	id := path.Base(request.URL.Path)
	deck := c.fetchDeck(id)
	if deck.LastShownIndex == 52 {
		message := struct {
			Remaining int `json:"remaining"`
		}{0}
		response.Header().Set("Content-Type", "application/json")
		json.NewEncoder(response).Encode(&message)
		return
	}
	nextCard := deck.Cards[deck.LastShownIndex]
	err := db.IncrementLastShown(c.DB, deck)
	if err != nil {
		panic(err)
	}
	message := struct {
		Card      models.Card `json:"card"`
		Remaining int         `json:"remaining"`
	}{nextCard, 52 - (deck.LastShownIndex + 1)}
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(&message)
}

func (c *Controller) Create(response http.ResponseWriter, request *http.Request) {
	newDeck := c.freshDeck()
	err := db.AddDeck(c.DB, newDeck)
	if err != nil {
		panic(err)
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(&newDeck.DeckID)
}

var suits = []string{"♥", "♣", "♦", "♠"}
var numbers = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "X", "J", "Q", "K"}

func (c *Controller) freshDeck() models.Deck {
	newDeck := models.Deck{}
	newDeck.LastShownIndex = 0
	newDeck.DeckID = c.newID()
	for _, suit := range suits {
		for _, number := range numbers {
			newCard := models.Card{Suit: suit, Number: number}
			newDeck.Cards = append(newDeck.Cards, newCard)
		}
		randomSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := len(newDeck.Cards) - 1; i > 0; i-- {
			j := randomSeed.Intn(i + 1)
			newDeck.Cards[i], newDeck.Cards[j] = newDeck.Cards[j], newDeck.Cards[i]
		}
	}
	return newDeck
}

func (c *Controller) newID() bson.ObjectId {
	for {
		new := bson.NewObjectId()
		unique, err := db.IsUnique(c.DB, new)
		if err != nil {
			panic(err)
		}
		if unique {
			return new
		}
	}
}
