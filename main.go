package main

import (
	"log"
	"os"
	"time"
)

// TODO: Get rid of this loop. Move notifications and filters to module code instead of app code.
func loop(config *AppConfig, tf *TurnipFinder) {
	islands := make([]Island, 0)

	for {
		newIslands := make([]Island, 0)
		pollingUsers := tf.PollingUsers()
		if len(pollingUsers) > 0 {
			newIslands = tf.PollSources()
		}

		if len(newIslands) > 0 {
			for _, island := range newIslands {

				islands = append(islands, island)
				log.Printf("[%d/%d] %s \tPrice: %d\tURL: %s\n", island.InQueue, island.MaxQueue, island.Name, island.TurnipPrice, island.URL)

				for _, user := range pollingUsers {
					if user.SellPrice > 0 && !FilterMinPrice(island, user.SellPrice) {
						break
					}
					if user.BuyPrice > 0 && !FilterMaxPrice(island, user.BuyPrice) {
						// TODO: Buying must also check for Daisy
						break
					}
					if len(user.ExcludePrices) > 0 && !FilterExcludePrices(island, user.ExcludePrices) {
						break
					}
					if user.MaxInQueue >= 0 && !FilterQueueSize(island, user.MaxInQueue) {
						break
					}

					err := tf.SendUserIsland(user, island)
					if err != nil {
						log.Println("Error sending island message")
						log.Fatal(err)
					}
				}
			}

		}

		time.Sleep(config.LoopInterval * time.Second)
	}
}

// TODO: Move main package to sub folder, for app, and change this package to turnipfinder.
func main() {
	// TODO: Check args & env variable. Store config and user state.
	config := NewConfig(os.Args[1])

	tf := New()
	tf.AddSource(NewTurnipExchangeSource())

	dg, err := DiscordConnect(config.DiscordBotToken)
	if err != nil {
		log.Fatal(err)
	}

	defer dg.Close()

	tf.RegisterDefaultCommands()

	tf.SendUserMessage = DiscordSendUserMessageWrapper(dg)

	dg.AddHandler(DiscordCreateMessageWrapper(tf))

	loop(config, tf)
}
