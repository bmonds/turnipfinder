package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func DiscordConnect(token string) (*discordgo.Session, error) {
	sess, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		return nil, err
	}

	err = sess.Open()
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func DiscordSendUserMessageWrapper(dg *discordgo.Session) func(user User, msg string) error {
	return func(user User, msg string) error {
		dgUser, err := dg.User(user.ID)
		if err != nil {
			return err
		}

		dgChannel, err := dg.UserChannelCreate(dgUser.ID)
		if err != nil {
			return err
		}

		_, err = dg.ChannelMessageSend(dgChannel.ID, msg)
		return err
	}
}

func DiscordCreateMessageWrapper(tf *TurnipFinder) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content[0:1] == "!" {
			fields := strings.Fields(m.Content)
			cmd := fields[0][1:]

			user, err := tf.User(m.Author.ID)
			if err != nil {
				user = tf.AddUserWithName(m.Author.ID, m.Author.Username)
			}

			reply := func(msg string) error {
				err := tf.SendUserMessage(user, msg)
				if err != nil {
					return err
				}

				return nil
			}

			commandInput := ChatCommandInput{
				Name:  cmd,
				Args:  strings.Join(fields[1:], " "), // Replace with Regex
				User:  user,
				Reply: reply,
			}

			err = tf.RunCommand(commandInput)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
