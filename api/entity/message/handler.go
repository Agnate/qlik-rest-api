package message

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type API struct {
	storage *MessageStorage
}

func New(db *sql.DB) *API {
	return &API{
		storage: NewMessageStorage(db),
	}
}

func (a *API) List(w http.ResponseWriter, r *http.Request) {
	msgs, err := a.storage.List()
	if err != nil {
		log.Fatalln(err)
		return
	}

	if len(msgs) <= 0 {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(msgs); err != nil {
		return
	}
}
