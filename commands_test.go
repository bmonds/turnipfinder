package main

import (
	"errors"
	"log"
	"regexp"
	"testing"
)

type mockedReply struct {
	Got []string
}

func (r *mockedReply) Add(msg string) {
	r.Got = append(r.Got, msg)
}

func mockReply(shouldError bool) (*mockedReply, func(string) error) {
	mock := mockedReply{
		Got: make([]string, 0),
	}

	reply := func(msg string) error {
		if shouldError {
			return errors.New("error")
		}

		mock.Add(msg)
		return nil
	}

	return &mock, reply
}

func TestCommandEcho(t *testing.T) {
	testTable := []struct {
		Input            ChatCommandInput
		Name             string
		ReplyShouldError bool
		ExpectedReplies  []string
		ExpectedError    bool
	}{
		{
			Name: "replies with the same message",
			Input: ChatCommandInput{
				Args: "foo bar",
			},
			ExpectedReplies: []string{"foo bar"},
			ExpectedError:   false,
		}, {
			Name: "returns the error from the reply",
			Input: ChatCommandInput{
				Args: "[error] failed", // Special string for tests to throw an error from reply()
			},
			ReplyShouldError: true,
			ExpectedReplies:  []string{},
			ExpectedError:    true,
		},
	}

	for _, tcase := range testTable {
		t.Run(tcase.Name, func(t *testing.T) {
			tf := New()
			mock, reply := mockReply(tcase.ReplyShouldError)
			tcase.Input.Reply = reply
			err := CommandEcho(tf, tcase.Input)

			if err != nil && !tcase.ExpectedError {
				t.Errorf("Expected nil to be returned by received an error")
			} else if err == nil && tcase.ExpectedError {
				t.Errorf("Expected an error to be returned but received nil")
			}

			if len(mock.Got) != len(tcase.ExpectedReplies) {
				t.Errorf("Expected %d replies but received %d", len(tcase.ExpectedReplies), len(mock.Got))
			}

			for idx, expected := range tcase.ExpectedReplies {
				if len(mock.Got) >= idx+1 {
					if mock.Got[idx] != expected {
						t.Errorf("Expected reply[%d] to be %q but received %q", idx, expected, mock.Got[idx])
					}
				} else {
					t.Errorf("Expected reply[%d] to be %q but it was not found.", idx, expected)
				}
			}
		})
	}
}

func TestCommandSell(t *testing.T) {
	userID := "foo"
	SetupUser := func(tf *TurnipFinder, user *User) {
		user.ID = userID
		tf.SetUser(*user)
	}
	testTable := []struct {
		Name                  string
		Input                 ChatCommandInput
		tf                    *TurnipFinder
		ReplyShouldError      bool
		ExpectedUserSellPrice int
		ExpectedUserPolling   bool
		ExpectedRepliesRegex  []*regexp.Regexp
		ExpectedError         bool
	}{
		{
			Name:                  "Shows usage on missing args",
			Input:                 ChatCommandInput{},
			ExpectedUserSellPrice: 0,
			ExpectedUserPolling:   false,
			ExpectedRepliesRegex:  []*regexp.Regexp{regexp.MustCompile(`Usage: .*`)},
			ExpectedError:         false,
		}, {
			Name: "Shows usage if args is a blank string",
			Input: ChatCommandInput{
				Args: "   ",
			},
			ExpectedUserSellPrice: 0,
			ExpectedUserPolling:   false,
			ExpectedRepliesRegex:  []*regexp.Regexp{regexp.MustCompile(`Usage: .*`)},
			ExpectedError:         false,
		}, {
			Name: "Shows usage if args is not a number",
			Input: ChatCommandInput{
				Args: "foo",
			},
			ExpectedUserSellPrice: 0,
			ExpectedUserPolling:   false,
			ExpectedRepliesRegex:  []*regexp.Regexp{regexp.MustCompile(`Usage: .*`)},
			ExpectedError:         false,
		}, {
			Name: "Sets the user's purchase price and enables polling",
			Input: ChatCommandInput{
				Args: "400",
			},
			ExpectedUserSellPrice: 400,
			ExpectedUserPolling:   true,
			ExpectedRepliesRegex:  []*regexp.Regexp{regexp.MustCompile(`.*will notify.*above 400.*`)},
			ExpectedError:         false,
		}, {
			Name: "Replaces the user's previous purchase price",
			Input: ChatCommandInput{
				Args: "200",
				User: User{
					SellPrice: 100,
				},
			},
			ExpectedUserSellPrice: 200,
			ExpectedUserPolling:   true,
			ExpectedRepliesRegex:  []*regexp.Regexp{regexp.MustCompile(`.*will notify.*above 200.*`)},
			ExpectedError:         false,
		}, {
			Name: "Returns the error from the reply",
			Input: ChatCommandInput{
				Args: "[error] failed", // Special string for tests to throw an error from reply()
			},
			ReplyShouldError:     true,
			ExpectedRepliesRegex: []*regexp.Regexp{},
			ExpectedError:        true,
		},
	}

	for _, tcase := range testTable {
		t.Run(tcase.Name, func(t *testing.T) {
			if tcase.tf == nil {
				tcase.tf = New()
			}

			SetupUser(tcase.tf, &tcase.Input.User)

			log.Println(tcase.tf.Users[userID])

			mock, reply := mockReply(tcase.ReplyShouldError)
			tcase.Input.Reply = reply
			err := CommandSell(tcase.tf, tcase.Input)
			log.Println(tcase.tf.Users[userID])
			log.Println(tcase.Input.User)

			user, _ := tcase.tf.User(userID)

			if user.SellPrice != tcase.ExpectedUserSellPrice {
				t.Errorf("Expected user's price filter to be %d but found %d", tcase.ExpectedUserSellPrice, user.SellPrice)
			}

			if user.Polling != tcase.ExpectedUserPolling {
				t.Errorf("Expected user's polling to be %t but found %t", tcase.ExpectedUserPolling, user.Polling)
			}

			if err != nil && !tcase.ExpectedError {
				t.Errorf("Expected nil to be returned by received an error")
			} else if err == nil && tcase.ExpectedError {
				t.Errorf("Expected an error to be returned but received nil")
			}

			if len(mock.Got) != len(tcase.ExpectedRepliesRegex) {
				t.Errorf("Expected %d replies but received %d", len(tcase.ExpectedRepliesRegex), len(mock.Got))
			}

			for idx, regex := range tcase.ExpectedRepliesRegex {
				if len(mock.Got) >= idx+1 {
					if !regex.MatchString(mock.Got[idx]) {
						t.Errorf("Expected reply[%d] to match /%s/ but received %q", idx, regex.String(), mock.Got[idx])
					}
				} else {
					t.Errorf("Expected reply[%d] to match /%s/ but it was not found.", idx, regex.String())
				}
			}
		})
	}
}
