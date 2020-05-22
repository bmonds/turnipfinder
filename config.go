package main

import (
	"time"
)

const defaultLoopInterval = 1

type AppConfig struct {
	DiscordBotToken string
	LoopInterval    time.Duration
}

func NewConfig(DiscordBotToken string) *AppConfig {
	return &AppConfig{
		DiscordBotToken: DiscordBotToken,
		LoopInterval:    defaultLoopInterval,
	}
}
