package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/agnate/qlikrestapi/internal/util"
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
	if err := a.outputList(users, w); err != nil {
		util.Status500APIError(w, errors.New("could not parse data to json"))
	}

	util.Status200APIOk(w)
}

func (a *API) outputList(users Users, w http.ResponseWriter) error {
	if len(users) <= 0 {
		util.APIJsonHeaders(w)
		fmt.Fprint(w, "[]")
		return nil
	}

	// TODO: Support JSON and XML by allowing user to pass optional
	// parameters to the API call to decide the format.
	if err := json.NewEncoder(w).Encode(users); err != nil {
		return err
	}

	util.APIJsonHeaders(w)
	return nil
}
