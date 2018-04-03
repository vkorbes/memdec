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

type controller struct {
	DB *mgo.Session
}

func NewController(db *mgo.Session) *controller {
	return &controller{
		DB: db,
	}
}

func (c *controller) ListAllDecks(response http.ResponseWriter, request *http.Request) {
	items, err := db.GetAllDecks(c.DB)
	if err != nil {
		panic(nil)
	}
	writeJSON(response, http.StatusOK, items)
}

func (c *controller) Info(response http.ResponseWriter, request *http.Request) {
	id := path.Base(request.URL.Path)
	deck := c.fetchDeck(id)
	writeJSON(response, http.StatusOK, deck)
}

func (c *controller) fetchDeck(id string) models.Deck {
	if !bson.IsObjectIdHex(id) {
		panic("invalid ID")
	}
	deck, err := db.GetDeck(c.DB, bson.ObjectIdHex(id))
	if err != nil {
		panic(err)
	}
	return deck
}

const cardsInDeck = 52

func (c *controller) NextCard(response http.ResponseWriter, request *http.Request) {
	id := path.Base(request.URL.Path)
	deck := c.fetchDeck(id)
	var nextCard interface{}
	var cardsRemaining int
	if deck.LastShownIndex != cardsInDeck {
		nextCard = deck.Cards[deck.LastShownIndex]
		err := db.IncrementLastShown(c.DB, deck)
		if err != nil {
			panic(err)
		}
		cardsRemaining = cardsInDeck - (deck.LastShownIndex + 1)
	}
	message := struct {
		Card      interface{} `json:"card,omitempty"`
		Remaining int         `json:"remaining"`
	}{nextCard, cardsRemaining}
	writeJSON(response, http.StatusOK, message)
}

func (c *controller) Create(response http.ResponseWriter, request *http.Request) {
	newDeck := c.freshDeck()
	err := db.AddDeck(c.DB, newDeck)
	if err != nil {
		panic(err)
	}
	writeJSON(response, http.StatusCreated, newDeck.DeckID)
}

var suits = []string{"♥", "♣", "♦", "♠"}
var numbers = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "X", "J", "Q", "K"}

func (c *controller) freshDeck() models.Deck {
	newDeck := models.Deck{DeckID: c.newID()}
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
		// rand.Shuffle
	}
	return newDeck
}

func (c *controller) newID() bson.ObjectId {
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

func writeJSON(response http.ResponseWriter, statusCode int, content interface{}) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	sb, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	_, err = response.Write(sb)
	if err != nil {
		panic(err)
	}
}
