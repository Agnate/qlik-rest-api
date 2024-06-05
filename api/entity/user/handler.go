package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/agnate/qlikrestapi/internal/apikey"
	"github.com/agnate/qlikrestapi/internal/util"
	"github.com/agnate/qlikrestapi/internal/util/baddata"
)

type API struct {
	storage *UserStorage
}

// Create a new Users API handler.
func New(db *sql.DB) *API {
	return &API{
		storage: NewUserStorage(db),
	}
}

// Retrieve a list of all Users.
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	users, err := a.storage.List()
	if err != nil {
		util.Status404NoAPIEndpoint(w, r, err)
		return
	}

	// Output list of users.
	if err := a.outputList(users, http.StatusOK, w); err != nil {
		util.Status500APIError(w, errors.New("could not parse data to json"))
	}
}

// Create and return a new User.
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	// Get data from POST body.
	userInput, err := a.getJsonBody(r)
	if err != nil {
		baddata.New400BadData(err).Render(w)
		return
	}

	// Validate and process user input.
	user, err := a.processUserInput(userInput)
	if err != nil {
		baddata.New400BadData(err).Render(w)
		return
	}

	// Create user.
	newUser, err := a.storage.Create(user)
	if err != nil {
		baddata.New400BadData(err).Render(w)
		return
	}

	// Overwrite API with raw key so the user can save it.
	newUser.APIKey = user.RawAPIKey

	// Output newly-created user.
	if err := a.outputSingle(newUser, http.StatusCreated, w); err != nil {
		util.Status500APIError(w, errors.New("could not parse data to json"))
	}
}

// Get user by their API key.
func (a *API) GetUserByAPIKey(rawAPIKey string) (*User, error) {
	// Hash the apiKey before searching database.
	bytes := apikey.HashAPIKey(rawAPIKey)
	hash := apikey.HashByteToString(bytes)

	// Look up User by hashed API key.
	user, err := a.storage.GetUserByAPIKey(hash)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (a *API) outputSingle(user *User, successHttpStatus int, w http.ResponseWriter) error {
	return a.outputList([]*User{user}, successHttpStatus, w)
}

func (a *API) outputList(users Users, successHttpStatus int, w http.ResponseWriter) error {
	if len(users) <= 0 {
		util.APIJsonHeaders(w)
		fmt.Fprint(w, "[]")
		return nil
	}

	// TODO: Support JSON and XML by allowing user to pass optional
	// parameters to the API call to decide the format.

	// Check JSON parsing for errors.
	jsonData, err := json.Marshal(users)
	if err != nil {
		return err
	}

	// Write success headers.
	util.APIJsonHeaders(w)
	w.WriteHeader(successHttpStatus)
	w.Write(jsonData)
	return nil
}

// Parse the JSON body of a request.
func (a *API) getJsonBody(r *http.Request) (*UserInput, error) {
	decoder := json.NewDecoder(r.Body)
	var userInput *UserInput
	err := decoder.Decode(&userInput)
	if err != nil {
		return nil, err
	}
	return userInput, nil
}

// Convert the a UserInput object to a User and fill in missing data.
// Only used for CREATE and UPDATE. Not needed for DELETE.
func (a *API) processUserInput(userInput *UserInput) (*User, error) {
	if len(userInput.Name) <= 0 {
		return nil, errors.New("you must provide a full_name")
	}

	if len(userInput.Email) <= 0 {
		return nil, errors.New("you must provide a valid email")
	}

	// TODO: Add email validation.

	// Create the base Message object for database storage.
	user := &User{
		Name:  userInput.Name,
		Email: userInput.Email,
	}

	// Generate an API key. We will return the raw key and store the hash.
	raw, hash := apikey.GenerateAPIKey()
	user.APIKey = apikey.HashByteToString(hash)
	user.RawAPIKey = raw

	return user, nil
}
