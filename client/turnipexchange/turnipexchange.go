package turnipexchange

import "net/http"

const defaultBaseURL = "https://api.turnip.exchange"
const defaultContentType = "application/json"

type Client struct {
	BaseURL     string
	ContentType string
}

type Island struct {
	Name        string
	Background  string
	Fruit       string
	TurnipPrice int
	MaxQueue    int
	TurnipCode  string `json:"turnipCode"`
	Hemisphere  string
	Watchlist   int
	Fee         int
	Islander    string
	Category    string
	IslandTime  string
	CreateTime  string // Change to date format
	Description string
	Queued      string // Parse as 2 fields.
	Patreon     int
	DiscordOnly int
	PatreonOnly int
	MessageID   string
	Thumbsupt   int
	Heart       int
	Poop        int
	Clown       int
	IslandScore float32
}

type ErrorWithResponse struct {
	Response http.Response
	Request  http.Request
}

type ErrorResponseNotSuccess struct {
	Success bool
	Message string
	ErrorWithResponse
}

type ErrorTryLater struct {
	ErrorWithResponse
}

func (e *ErrorResponseNotSuccess) Error() string {
	return e.Message
}

func (e *ErrorTryLater) Error() string {
	return "Try Later"
}

func New() *Client {
	return &Client{
		BaseURL:     defaultBaseURL,
		ContentType: defaultContentType,
	}
}
