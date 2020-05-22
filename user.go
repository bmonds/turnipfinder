package main

type User struct {
	ID            string
	Name          string
	Polling       bool
	SellPrice     int
	BuyPrice      int
	ExcludePrices []int
	MaxInQueue    int
}

type ErrorUserNotFound struct{}

func (e *ErrorUserNotFound) Error() string {
	return "User was not found"
}

func (tf *TurnipFinder) AddUserWithName(ID string, Name string) User {
	tf.Users[ID] = User{
		ID:            ID,
		Name:          Name,
		SellPrice:     0,
		BuyPrice:      0,
		ExcludePrices: []int{666},
		MaxInQueue:    -1,
		Polling:       false,
	}

	return tf.Users[ID]
}

func (tf *TurnipFinder) AddUser(ID string) User {
	return tf.AddUserWithName(ID, ID)
}

func (tf *TurnipFinder) SetUser(user User) {
	tf.Users[user.ID] = user
}

func (tf *TurnipFinder) User(ID string) (User, error) {
	if user, ok := tf.Users[ID]; ok {
		return user, nil
	}

	return User{}, &ErrorUserNotFound{}
}
