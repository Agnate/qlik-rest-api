package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/agnate/qlikrestapi/internal/util"
)

type API struct {
	storage *UserStorage
}

func New(db *sql.DB) *API {
	return &API{
		storage: NewUserStorage(db),
	}
}

func (a *API) List(w http.ResponseWriter, r *http.Request) {
	users, err := a.storage.List()
	if err != nil {
		util.NoAPIEndpoint(w, err)
		return
	}

	if len(users) <= 0 {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		return
	}
}
