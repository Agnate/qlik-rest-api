package baddata

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/agnate/qlikrestapi/internal/util"
)

type BadData struct {
	ErrorMsg string
	err      error
}

// Create a BadData entry to be displayed to the user in the API body.
func New400BadData(err error) *BadData {
	return &BadData{
		ErrorMsg: err.Error(),
		err:      err,
	}
}

// Used when validation or data saving fails for an endpoint and we want a consistent
// output displayed to our users. Errors will be logged.
func (bd *BadData) Render(w http.ResponseWriter) {
	// TODO: Add logging for invalid endpoints in case we need to monitor spammers.
	log.Println(bd.err)
	http.Error(w, util.NewHttpStatusMsg(http.StatusBadRequest), http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(bd); err != nil {
		return
	}
}
