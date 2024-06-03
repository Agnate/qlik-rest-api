package message

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	myCtx "github.com/agnate/qlikrestapi/internal/context"
	"github.com/agnate/qlikrestapi/internal/util"
	"github.com/agnate/qlikrestapi/internal/util/baddata"
	"github.com/agnate/qlikrestapi/internal/validation"
	"github.com/google/uuid"
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
	// List out data from storage.
	msgs, err := a.storage.List()
	if err != nil {
		util.NoAPIEndpoint(w, err)
		return
	}
	a.outputList(msgs, w)
}

func (a *API) ListByUUID(w http.ResponseWriter, r *http.Request) {
	// Validate route data from context.
	validUUID, err := a.validateUUID(r)
	if err != nil {
		util.NoAPIEndpoint(w, err)
		return
	}

	// List out data from storage.
	msgs, err := a.storage.ListByUUID(validUUID.Parsed)
	if err != nil {
		util.NoAPIEndpoint(w, err)
		return
	}
	a.outputList(msgs, w)
}

func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	// Validate route data from context.
	validUUID, validCreateDate, err := a.validatePrimaryKey(r)
	if err != nil {
		util.NoAPIEndpoint(w, err)
		return
	}

	// Read in data from storage.
	msg, err := a.storage.Read(validUUID.Parsed, validCreateDate.Parsed)
	if err != nil {
		util.NoAPIEndpoint(w, err)
		return
	}
	a.outputSingle(msg, w)
}

func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	// Validate route data from context.
	validUUID, err := a.validateUUID(r)
	if err != nil {
		baddata.New(err).Render(w)
		return
	}

	// Get data from POST body.
	msgInput, err := a.getJsonBody(r)
	if err != nil {
		baddata.New(err).Render(w)
		return
	}

	// Validate and process message input.
	msg, err := a.processMessageInput(msgInput, validUUID.Parsed, time.Time{})
	if err != nil {
		baddata.New(err).Render(w)
		return
	}

	// Create message.
	newMsg, err := a.storage.Create(msg)
	if err != nil {
		baddata.New(err).Render(w)
		return
	}

	// Output newly-created message.
	a.outputSingle(newMsg, w)
}

func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	// Validate route data from context.
	validUUID, validCreateDate, err := a.validatePrimaryKey(r)
	if err != nil {
		util.NoAPIEndpoint(w, err)
		return
	}

	// Load existing Message so we can check concurrency.
	existingMsg, err := a.storage.Read(validUUID.Parsed, validCreateDate.Parsed)
	if err != nil || existingMsg == nil {
		util.NoAPIEndpoint(w, err)
		return
	}

	// Get data from POST body.
	msgInput, err := a.getJsonBody(r)
	if err != nil {
		baddata.New(err).Render(w)
		return
	}

	// Check concurrency before processing.
	if !a.isConcurrent(existingMsg, msgInput) {
		a.getConcurrentBadData(msgInput).Render(w)
		return
	}

	// Validate and process message input.
	msg, err := a.processMessageInput(msgInput, validUUID.Parsed, validCreateDate.Parsed)
	if err != nil {
		baddata.New(err).Render(w)
		return
	}

	// Update message.
	updatedMsg, err := a.storage.Update(msg)
	if err != nil {
		baddata.New(err).Render(w)
		return
	}

	// Output updated message.
	a.outputSingle(updatedMsg, w)
}

func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	// Validate route data from context.
	validUUID, validCreateDate, err := a.validatePrimaryKey(r)
	if err != nil {
		util.NoAPIEndpoint(w, err)
		return
	}

	// Load existing Message so we can check concurrency.
	existingMsg, err := a.storage.Read(validUUID.Parsed, validCreateDate.Parsed)
	if err != nil || existingMsg == nil {
		util.NoAPIEndpoint(w, err)
		return
	}

	// Get data from POST body.
	msgInput, err := a.getJsonBody(r)
	if err != nil {
		baddata.New(err).Render(w)
		return
	}

	// Check concurrency before processing.
	if !a.isConcurrent(existingMsg, msgInput) {
		a.getConcurrentBadData(msgInput).Render(w)
		return
	}

	// Note: No need to process MessageInput as we will use the existing message.

	// Delete message.
	deletedMsg, err := a.storage.Delete(existingMsg)
	if err != nil {
		baddata.New(err).Render(w)
		return
	}

	// Output deleted message.
	a.outputSingle(deletedMsg, w)
}

func (a *API) outputSingle(msg *Message, w http.ResponseWriter) {
	a.outputList([]*Message{msg}, w)
}

func (a *API) outputList(msgs Messages, w http.ResponseWriter) {
	if len(msgs) <= 0 {
		fmt.Fprint(w, "[]")
		return
	}

	// TODO: Support JSON and XML by allowing user to pass optional
	// parameters to the API call to decide the format.
	if err := json.NewEncoder(w).Encode(msgs); err != nil {
		return
	}
}

// Parse the JSON body of a request.
func (a *API) getJsonBody(r *http.Request) (*MessageInput, error) {
	decoder := json.NewDecoder(r.Body)
	var msg *MessageInput
	err := decoder.Decode(&msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// Convert the a MessageInput object to a Message and fill in missing data.
// Only used for CREATE and UPDATE. Not needed for DELETE.
func (a *API) processMessageInput(msgInput *MessageInput, uuid uuid.UUID, createDate time.Time) (*Message, error) {
	isUpdating := !createDate.IsZero()

	if len(msgInput.Message) <= 0 {
		return nil, errors.New("no message provided or is empty")
	}

	if isUpdating && msgInput.LastUpdated.IsZero() {
		return nil, errors.New("you must provide the most recent last_updated_date to modify this message")
	}

	// Create the base Message object for database storage.
	msg := &Message{
		Message: msgInput.Message,
	}
	msg.UUID = uuid
	msg.LastUpdatedBy = uuid
	msg.Palindrome = util.IsPalindrome(msg.Message)

	// If there is a createDate, it means we are updating an existing message,
	// so we need to set the CreateDate and LastUpdated fields as well.
	if isUpdating {
		msg.CreateDate = createDate
		msg.LastUpdated = msgInput.LastUpdated
	}
	return msg, nil
}

func (a *API) isConcurrent(existingMsg *Message, msgInput *MessageInput) bool {
	// Note: We also need to check concurreny as part of our database query.
	return existingMsg.LastUpdated.Equal(msgInput.LastUpdated)
}

func (a *API) getConcurrentBadData(msgInput *MessageInput) *baddata.BadData {
	err := errors.New("this message has been updated by someone else - please resubmit with most recent last_updated_date")
	if msgInput.LastUpdated.IsZero() {
		err = errors.New("you must provide the most recent last_updated_date to modify this message")
	}
	return baddata.New(err)
}

func (a *API) validateUUID(r *http.Request) (*validation.RuleUUID, error) {
	// TODO: Improve slug management so handler doesn't need to know the index
	uuidSlugIndex := 0

	// Get route data from context.
	rawUuid := myCtx.GetSlug(r.Context(), uuidSlugIndex)

	// Validate API slugs.
	return validation.NewRuleUUID(rawUuid)
}

func (a *API) validatePrimaryKey(r *http.Request) (*validation.RuleUUID, *validation.RuleTime, error) {
	// TODO: Improve slug management so handler doesn't need to know the index
	uuidSlugIndex := 0
	createDateSlugIndex := 1

	// Get route data from context.
	rawUuid := myCtx.GetSlug(r.Context(), uuidSlugIndex)
	rawCreateDate := myCtx.GetSlug(r.Context(), createDateSlugIndex)

	// Validate API slugs.
	validUUID, err := validation.NewRuleUUID(rawUuid)
	if err != nil {
		return nil, nil, err
	}

	validCreateDate, err := validation.NewRuleTime(time.RFC3339Nano, rawCreateDate)
	if err != nil {
		return validUUID, nil, err
	}

	return validUUID, validCreateDate, err
}
