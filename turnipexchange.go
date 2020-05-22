package main

import (
	"fmt"
	"github.com/bmonds/turnipfinder/client/turnipexchange"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type TurnipExchangeSource struct {
	client        *turnipexchange.Client
	lastRateLimit TurnipExchangeRateLimit
}

type TurnipExchangeRateLimit struct {
	Limit     int
	Remaining int
	Reset     int64
	Next      int64
}

func NewTurnipExchangeSource() *TurnipExchangeSource {
	return &TurnipExchangeSource{
		client: turnipexchange.New(),
	}
}

func (t *TurnipExchangeSource) ToIsland(island turnipexchange.Island) Island {
	inQueue := -1
	regex := regexp.MustCompile(`^(\d+)\/(\d+)$`)
	match := regex.FindStringSubmatch(island.Queued)
	if len(match) == 3 {
		val, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatal(err)
		}

		inQueue = val
	}

	return Island{
		Name:        island.Name,
		TurnipPrice: island.TurnipPrice,
		MaxQueue:    island.MaxQueue,
		URL:         fmt.Sprintf("https://turnip.exchange/island/%s", island.TurnipCode),
		Fee:         island.Fee,
		Islander:    island.Islander,
		Category:    island.Category,
		CreateTime:  island.CreateTime,
		Description: island.Description,
		InQueue:     inQueue,
	}
}

func (t *TurnipExchangeSource) SetTurnipExchangeRateLimit(headers http.Header) {
	limit, err := strconv.Atoi(headers.Get("X-Ratelimit-Limit"))
	if err != nil {
		log.Fatal(err)
	}

	remaining, err := strconv.Atoi(headers.Get("X-Ratelimit-Remaining"))
	if err != nil {
		log.Fatal(err)
	}

	reset, err := strconv.ParseInt(headers.Get("X-Ratelimit-Reset"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	next := reset
	if remaining > 0 {
		now := time.Now()
		epoch := now.Unix()
		diff := reset - epoch
		secondsPer := diff / int64(remaining)
		next = epoch + secondsPer
	}

	t.lastRateLimit = TurnipExchangeRateLimit{
		Limit:     limit,
		Remaining: remaining,
		Reset:     reset,
		Next:      next,
	}
}

func (t *TurnipExchangeSource) TurnipExchangeIsRateLimited() bool {
	now := time.Now()
	epoch := now.Unix()

	if t.lastRateLimit.Remaining == 0 && epoch < t.lastRateLimit.Reset {
		log.Println("Rate Limited")
		return true
	}
	if t.lastRateLimit.Next > epoch {
		return true
	}

	return false
}

func (t *TurnipExchangeSource) Run() []Island {
	islands := make([]Island, 0)

	if t.TurnipExchangeIsRateLimited() {
		return islands
	}

	teIslands, resp, err := t.client.Islands("neither", "turnips", 0)
	if err != nil {
		if _, ok := err.(*turnipexchange.ErrorTryLater); !ok {
			log.Fatal(err)
		}
	}

	if resp != nil {
		t.SetTurnipExchangeRateLimit(resp.Header)
	}

	for _, island := range teIslands {
		islands = append(islands, t.ToIsland(island))
	}

	return islands
}
