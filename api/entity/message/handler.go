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

// Create a new Messages API handler.
func New(db *sql.DB) *API {
	return &API{
		storage: NewMessageStorage(db),
	}
}

// Retrieve a list of all Messages.
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	// List out data from storage.
	msgs, err := a.storage.List()
	if err != nil {
		util.Status404NoAPIEndpoint(w, r, err)
		return
	}

	// Output list of messages.
	if err := a.outputList(msgs, w); err != nil {
		util.Status500APIError(w, errors.New("could not parse data to json"))
	}

	util.Status200APIOk(w)
}

// Retrieve list of all Messages for specific User.
func (a *API) ListByUUID(w http.ResponseWriter, r *http.Request) {
	// Validate route data from context.
	validUUID, err := a.validateUUID(r)
	if err != nil {
		util.Status404NoAPIEndpoint(w, r, err)
		return
	}

	// List out data from storage.
	msgs, err := a.storage.ListByUUID(validUUID.Parsed)
	if err != nil {
		util.Status404NoAPIEndpoint(w, r, err)
		return
	}

	// Output list of messages.
	if err := a.outputList(msgs, w); err != nil {
		util.Status500APIError(w, errors.New("could not parse data to json"))
	}

	util.Status200APIOk(w)
}

// Return a specific Message based on primary key (UUID, CreateDate).
func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	// Validate route data from context.
	validUUID, validCreateDate, err := a.validatePrimaryKey(r)
	if err != nil {
		util.Status404NoAPIEndpoint(w, r, err)
		return
	}

	// Read in data from storage.
	msg, err := a.storage.Read(validUUID.Parsed, validCreateDate.Parsed)
	if err != nil {
		util.Status404NoAPIEndpoint(w, r, err)
		return
	}

	// Output the message.
	if err := a.outputSingle(msg, w); err != nil {
		util.Status500APIError(w, errors.New("could not parse data to json"))
	}

	util.Status200APIOk(w)
}

// Create and return a new Message.
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	// Validate route data from context.
	validUUID, err := a.validateUUID(r)
	if err != nil {
		baddata.New400BadData(err).Render(w)
		return
	}

	// Get data from POST body.
	msgInput, err := a.getJsonBody(r)
	if err != nil {
		baddata.New400BadData(err).Render(w)
		return
	}

	// Validate and process message input.
	msg, err := a.processMessageInput(msgInput, validUUID.Parsed, time.Time{})
	if err != nil {
		baddata.New400BadData(err).Render(w)
		return
	}

	// Create message.
	newMsg, err := a.storage.Create(msg)
	if err != nil {
		baddata.New400BadData(err).Render(w)
		return
	}

	// Output newly-created message.
	if err := a.outputSingle(newMsg, w); err != nil {
		util.Status500APIError(w, errors.New("could not parse data to json"))
	}

	util.Status201APICreate(w)
}

// Update and return an existing Message.
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	// Validate route data from context.
	validUUID, validCreateDate, err := a.validatePrimaryKey(r)
	if err != nil {
		util.Status404NoAPIEndpoint(w, r, err)
		return
	}

	// Load existing Message so we can check concurrency.
	existingMsg, err := a.storage.Read(validUUID.Parsed, validCreateDate.Parsed)
	if err != nil || existingMsg == nil {
		util.Status404NoAPIEndpoint(w, r, err)
		return
	}

	// Get data from POST body.
	msgInput, err := a.getJsonBody(r)
	if err != nil {
		baddata.New400BadData(err).Render(w)
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
		baddata.New400BadData(err).Render(w)
		return
	}

	// Update message.
	updatedMsg, err := a.storage.Update(msg)
	if err != nil {
		baddata.New400BadData(err).Render(w)
		return
	}

	// Output updated message.
	if err := a.outputSingle(updatedMsg, w); err != nil {
		util.Status500APIError(w, errors.New("could not parse data to json"))
	}

	util.Status200APIOk(w)
}

// Delete an existing Message.
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	// Validate route data from context.
	validUUID, validCreateDate, err := a.validatePrimaryKey(r)
	if err != nil {
		util.Status404NoAPIEndpoint(w, r, err)
		return
	}

	// Load existing Message so we can check concurrency.
	existingMsg, err := a.storage.Read(validUUID.Parsed, validCreateDate.Parsed)
	if err != nil || existingMsg == nil {
		util.Status404NoAPIEndpoint(w, r, err)
		return
	}

	// Get data from POST body.
	msgInput, err := a.getJsonBody(r)
	if err != nil {
		baddata.New400BadData(err).Render(w)
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
		baddata.New400BadData(err).Render(w)
		return
	}

	// Output deleted message.
	if err := a.outputSingle(deletedMsg, w); err != nil {
		util.Status500APIError(w, errors.New("could not parse data to json"))
	}

	util.Status200APIOk(w)
}

func (a *API) outputSingle(msg *Message, w http.ResponseWriter) error {
	return a.outputList([]*Message{msg}, w)
}

func (a *API) outputList(msgs Messages, w http.ResponseWriter) error {
	if len(msgs) <= 0 {
		util.APIJsonHeaders(w)
		fmt.Fprint(w, "[]")
		return nil
	}

	// TODO: Support JSON and XML by allowing user to pass optional
	// parameters to the API call to decide the format.
	if err := json.NewEncoder(w).Encode(msgs); err != nil {
		return err
	}

	util.APIJsonHeaders(w)
	return nil
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
	return baddata.New400BadData(err)
}

func (a *API) validateUUID(r *http.Request) (*validation.RuleUUID, error) {
	// TODO: Improve slug management so handler doesn't need to know the index
	uuidSlugIndex := 0

	// Get route data from context.
	rawUuid, err := myCtx.GetSlug(r.Context(), uuidSlugIndex)
	if err != nil {
		return nil, err
	}

	// Validate API slugs.
	return validation.NewRuleUUID(rawUuid)
}

func (a *API) validatePrimaryKey(r *http.Request) (*validation.RuleUUID, *validation.RuleTime, error) {
	// TODO: Improve slug management so handler doesn't need to know the index
	uuidSlugIndex := 0
	createDateSlugIndex := 1

	// Get route data from context.
	rawUuid, err := myCtx.GetSlug(r.Context(), uuidSlugIndex)
	if err != nil {
		return nil, nil, err
	}

	rawCreateDate, err := myCtx.GetSlug(r.Context(), createDateSlugIndex)
	if err != nil {
		return nil, nil, err
	}

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
