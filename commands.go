package main

import (
	"fmt"
	"strconv"
	"strings"
)

type ChatCommand func(tf *TurnipFinder, input ChatCommandInput) error

type ChatCommandInput struct {
	Name  string
	Args  string
	User  User
	Reply func(string) error
}

func FormatCommandName(name string) string {
	return strings.Trim(strings.ToLower(name), " ")
}

func (tf *TurnipFinder) AddCommand(name string, f ChatCommand) {
	mapName := FormatCommandName(name)

	tf.commands[mapName] = f
}

func (tf *TurnipFinder) RegisterDefaultCommands() {
	tf.AddCommand("help", CommandHelp)
	tf.AddCommand("echo", CommandEcho)
	tf.AddCommand("sell", CommandSell)
	tf.AddCommand("buy", CommandBuy)
	tf.AddCommand("maxqueue", CommandMaxQueue)
	tf.AddCommand("stop", CommandStop)
	tf.AddCommand("status", CommandStatus)
}

func (tf *TurnipFinder) GetCommand(name string) ChatCommand {
	mapName := FormatCommandName(name)
	if tf.commands[mapName] == nil {
		return nil
	}

	return tf.commands[mapName]
}

func (tf *TurnipFinder) RunCommand(input ChatCommandInput) error {
	cmd := tf.GetCommand(input.Name)
	err := cmd(tf, input)
	if err != nil {
		return err
	}

	return nil
}

func CommandEcho(tf *TurnipFinder, input ChatCommandInput) error {
	return input.Reply(input.Args)
}

func CommandHelp(tf *TurnipFinder, input ChatCommandInput) error {
	arrCommands := make([]string, 0)
	for name := range tf.commands {
		arrCommands = append(arrCommands, name)
	}
	return input.Reply(fmt.Sprintf("Commands: %s", strings.Join(arrCommands, ", ")))
}

func CommandSell(tf *TurnipFinder, input ChatCommandInput) error {
	price, err := strconv.Atoi(input.Args)
	if err != nil {
		return input.Reply("Usage: !sell [minPrice]")
	}

	if price < tf.MinTurnipPriceAllowed || price > tf.MaxTurnipPriceAllowed {
		err := input.Reply(fmt.Sprintf("Sell price must be between %d and %d", tf.MinTurnipPriceAllowed, tf.MaxTurnipPriceAllowed))
		if err != nil {
			return err
		}
	}

	input.User.SellPrice = price
	input.User.Polling = true
	tf.SetUser(input.User)

	return input.Reply(fmt.Sprintf("I will notify you about islands buying turnips above %d", price))
}

func CommandBuy(tf *TurnipFinder, input ChatCommandInput) error {
	price, err := strconv.Atoi(input.Args)
	if err != nil {
		return input.Reply("Usage: !buy [minPrice]")
	}

	if price < tf.MinTurnipPriceAllowed || price > tf.MaxTurnipPriceAllowed {
		return input.Reply(fmt.Sprintf("Sell price must be between %d and %d", tf.MinTurnipPriceAllowed, tf.MaxTurnipPriceAllowed))
	}

	input.User.BuyPrice = price
	input.User.Polling = true
	tf.SetUser(input.User)

	return input.Reply(fmt.Sprintf("I will notify you about islands selling turnips below %d", price))
}

func CommandMaxQueue(tf *TurnipFinder, input ChatCommandInput) error {
	maxInQueue, err := strconv.Atoi(input.Args)
	if err != nil {
		return input.Reply("Usage: !maxqueue [maxUsersInQueue]")
	}

	input.User.MaxInQueue = maxInQueue

	tf.SetUser(input.User)
	return input.Reply(fmt.Sprintf("I will only send items that have %d users in the queue or less", maxInQueue))
}

func CommandStop(tf *TurnipFinder, input ChatCommandInput) error {
	input.User.Polling = false

	tf.SetUser(input.User)
	return input.Reply("You have stopped looking for an island.")
}

func CommandStatus(tf *TurnipFinder, input ChatCommandInput) error {
	msgPolling := "You are not currently looking for islands"
	if input.User.Polling {
		msgPolling = "You are currently looking for islands"

		if input.User.SellPrice > 0 {
			msgPolling += fmt.Sprintf(" with a turnip price over %d", input.User.SellPrice)
		} else if input.User.BuyPrice > 0 {
			msgPolling += fmt.Sprintf(" with a turnip price under %d", input.User.BuyPrice)
		}
	}

	return input.Reply(msgPolling)
}
