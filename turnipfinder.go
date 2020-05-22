package main

import "fmt"

const (
	defaultMinTurnipPriceAllowed = 15
	defaultMaxTurnipPriceAllowed = 800
)

type TurnipFinder struct {
	Config                TurnipFinderConfig
	Users                 map[string]User
	Islands               map[string]Island
	Sources               []IslandSource
	MinTurnipPriceAllowed int
	MaxTurnipPriceAllowed int
	SendUserMessage       SendUserMessage
	commands              map[string]ChatCommand
}

type TurnipFinderConfig struct {
}

type IslandSource interface {
	Run() []Island
}

type Island struct {
	ID          string
	Name        string
	TurnipPrice int
	MaxQueue    int
	URL         string
	Fee         int
	Islander    string
	Category    string
	IslandTime  string
	CreateTime  string // Change to date format
	Description string
	InQueue     int
}

type SendUserMessage func(user User, message string) error

func New() *TurnipFinder {
	return &TurnipFinder{
		MinTurnipPriceAllowed: defaultMinTurnipPriceAllowed,
		MaxTurnipPriceAllowed: defaultMaxTurnipPriceAllowed,
		Sources:               make([]IslandSource, 0),
		Users:                 make(map[string]User),
		Islands:               make(map[string]Island),
		commands:              make(map[string]ChatCommand),
	}
}

func (tf *TurnipFinder) PollSources() []Island {
	// TODO: Move to goroutines
	newIslands := make([]Island, 0)

	for idx := range tf.Sources {
		islands := tf.Sources[idx].Run()
		for _, island := range islands {
			if _, ok := tf.Islands[island.ID]; !ok {
				newIslands = append(newIslands, island)
			}
			tf.Islands[island.ID] = island
		}

	}

	return newIslands
}

func (tf *TurnipFinder) AddSource(source IslandSource) {
	tf.Sources = append(tf.Sources, source)
}

func (tf *TurnipFinder) PollingUsers() []User {
	users := make([]User, 0)

	for _, user := range tf.Users {
		if user.Polling {
			users = append(users, user)
		}
	}

	return users
}

func (tf *TurnipFinder) SendUserIsland(user User, island Island) error {
	msg := fmt.Sprintf("[%d/%d] %s \tPrice: %d\nURL: %s\nFee: %d\n%s\n", island.InQueue, island.MaxQueue, island.Name, island.TurnipPrice, island.URL, island.Fee, island.Description)
	err := tf.SendUserMessage(user, msg)

	return err
}
