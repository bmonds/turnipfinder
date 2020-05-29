package main

import "testing"

func TestAddUserWithName(t *testing.T) {
	testTable := []struct {
		Name          string
		tf            *TurnipFinder
		UserID        string
		UserName      string
		ExpectedUsers []User
	}{
		{
			Name:     "Create first user",
			UserID:   "foo",
			UserName: "bar",
			ExpectedUsers: []User{
				{ID: "foo", Name: "bar"},
			},
		},
		{
			Name: "Add additional users",
			tf: func() *TurnipFinder {
				newTf := New()
				newTf.Users["bar"] = User{
					ID:   "bar",
					Name: "bar user",
				}

				return newTf
			}(),
			UserID:   "foo",
			UserName: "bar",
			ExpectedUsers: []User{
				{ID: "foo", Name: "bar"},
				{ID: "bar", Name: "bar user"},
			},
		},
	}

	for _, tcase := range testTable {
		t.Run(tcase.Name, func(t *testing.T) {
			if tcase.tf == nil {
				tcase.tf = New()
			}

			tcase.tf.AddUserWithName(tcase.UserID, tcase.UserName)

			if len(tcase.tf.Users) != len(tcase.ExpectedUsers) {
				t.Errorf("Expected to find %d users but found %d", len(tcase.tf.Users), len(tcase.ExpectedUsers))
			}

			for _, expectedUser := range tcase.ExpectedUsers {
				user, ok := tcase.tf.Users[expectedUser.ID]
				if !ok {
					t.Errorf("Could not find user with ID %s", expectedUser.ID)
				} else if user.ID != expectedUser.ID {
					t.Errorf("Expected user at index %s to have ID %s but found %s", expectedUser.ID, expectedUser.ID, user.ID)
				} else if user.Name != expectedUser.Name {
					t.Errorf("Expected user with ID %s to have name %s but found %s", expectedUser.ID, expectedUser.Name, user.Name)
				}
			}
		})
	}
}

func TestAddUser(t *testing.T) {
	// AddUser() wraps AddUserWithName() to set a default name value
	testTable := []struct {
		Name          string
		tf            *TurnipFinder
		UserID        string
		ExpectedUsers []User
	}{
		{
			Name:   "Create user using the ID for the name",
			UserID: "foo",
			ExpectedUsers: []User{
				{ID: "foo", Name: "foo"},
			},
		},
	}

	for _, tcase := range testTable {
		t.Run(tcase.Name, func(t *testing.T) {
			if tcase.tf == nil {
				tcase.tf = New()
			}

			tcase.tf.AddUser(tcase.UserID)

			if len(tcase.tf.Users) != len(tcase.ExpectedUsers) {
				t.Errorf("Expected to find %d users but found %d", len(tcase.tf.Users), len(tcase.ExpectedUsers))
			}

			for _, expectedUser := range tcase.ExpectedUsers {
				user, ok := tcase.tf.Users[expectedUser.ID]
				if !ok {
					t.Errorf("Could not find user with ID %s", expectedUser.ID)
				} else if user.ID != expectedUser.ID {
					t.Errorf("Expected user at index %s to have ID %s but found %s", expectedUser.ID, expectedUser.ID, user.ID)
				} else if user.Name != expectedUser.Name {
					t.Errorf("Expected user with ID %s to have name %s but found %s", expectedUser.ID, expectedUser.Name, user.Name)
				}
			}
		})
	}
}

func TestSetUser(t *testing.T) {
	testTable := []struct {
		Name          string
		tf            *TurnipFinder
		User          User
		ExpectedUsers []User
	}{
		{
			Name: "Create user with the provided User object",
			User: User{ID: "foo", Name: "bar", SellPrice: 123, Polling: true},
			ExpectedUsers: []User{
				{ID: "foo", Name: "bar", SellPrice: 123, Polling: true},
			},
		},
	}

	for _, tcase := range testTable {
		t.Run(tcase.Name, func(t *testing.T) {
			if tcase.tf == nil {
				tcase.tf = New()
			}

			tcase.tf.SetUser(tcase.User)

			if len(tcase.tf.Users) != len(tcase.ExpectedUsers) {
				t.Errorf("Expected to find %d users but found %d", len(tcase.tf.Users), len(tcase.ExpectedUsers))
			}

			for _, expectedUser := range tcase.ExpectedUsers {
				user, ok := tcase.tf.Users[expectedUser.ID]
				if !ok {
					t.Errorf("Could not find user with ID %s", expectedUser.ID)
				} else if user.ID != expectedUser.ID {
					t.Errorf("Expected user at index %s to have ID %s but found %s", expectedUser.ID, expectedUser.ID, user.ID)
				} else if user.Name != expectedUser.Name {
					t.Errorf("Expected user with ID %s to have name %s but found %s", expectedUser.ID, expectedUser.Name, user.Name)
				} else if user.SellPrice != expectedUser.SellPrice {
					t.Errorf("Expected user with ID %s to have sell price of %d but found %d", expectedUser.ID, expectedUser.SellPrice, user.SellPrice)
				} else if user.Polling != expectedUser.Polling {
					t.Errorf("Expected user with ID %s to have polling value of %t but found %t", expectedUser.ID, expectedUser.Polling, user.Polling)
				}
			}
		})
	}
}

func TestUser(t *testing.T) {
	createTestUser := func() User {
		return User{ID: "foo", Name: "bar", SellPrice: 123, Polling: true}
	}
	testTable := []struct {
		Name          string
		tf            *TurnipFinder
		UserID        string
		ExpectedUser  User
		ExpectedError bool
	}{
		{
			Name: "Returns the requested User object",
			tf: func() *TurnipFinder {
				newTf := New()
				newTf.Users["foo"] = createTestUser()

				return newTf
			}(),
			UserID:       "foo",
			ExpectedUser: User{ID: "foo", Name: "bar", SellPrice: 123, Polling: true},
		},
		{
			Name: "Returns an error if the user was not found",
			tf: func() *TurnipFinder {
				newTf := New()
				newTf.Users["foo"] = createTestUser()

				return newTf
			}(),
			UserID:        "bar",
			ExpectedError: true,
		},
	}

	for _, tcase := range testTable {
		t.Run(tcase.Name, func(t *testing.T) {
			if tcase.tf == nil {
				tcase.tf = New()
			}

			user, err := tcase.tf.User(tcase.UserID)

			if tcase.ExpectedError {
				if err == nil {
					t.Error("Expected an error to be returned")
				}
			} else if err != nil {
				t.Error("Expected error to be nil")
			} else {
				if user.ID != tcase.ExpectedUser.ID {
					t.Errorf("Expected user to have ID %q but found %q", tcase.ExpectedUser.ID, user.ID)
				}
				if user.Name != tcase.ExpectedUser.Name {
					t.Errorf("Expected user to have name %q but found %q", tcase.ExpectedUser.Name, user.Name)
				}
				if user.SellPrice != tcase.ExpectedUser.SellPrice {
					t.Errorf("Expected user to have sell price of %d but found %d", tcase.ExpectedUser.SellPrice, user.SellPrice)
				}
				if user.Polling != tcase.ExpectedUser.Polling {
					t.Errorf("Expected user to have polling value of %t but found %t", tcase.ExpectedUser.Polling, user.Polling)
				}
			}
		})
	}
}
