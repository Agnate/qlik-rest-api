package message

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var listSelectQuery string = "^SELECT (.+) FROM messages WHERE logical_delete = \\$1$"

func TestMessageListHasResult(t *testing.T) {
	// Mock the database.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set up some mocked return data.
	wantRows := 1
	rows := getMessageRows(wantRows)

	// We expect a query to be run.
	mock.ExpectQuery(listSelectQuery).WithArgs(false).WillReturnRows(rows)

	// Instantiate storage.
	storage := NewMessageStorage(db)

	// Run and validate.
	if msgs, err := storage.List(); err != nil || len(msgs) != wantRows {
		t.Errorf("error was not expected while listing messages: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMessageListNoResults(t *testing.T) {
	// Mock the database.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set up some mocked return data.
	wantRows := 0
	rows := getMessageRows(wantRows)

	// We expect a query to be run.
	mock.ExpectQuery(listSelectQuery).WithArgs(false).WillReturnRows(rows)

	// Instantiate storage.
	storage := NewMessageStorage(db)

	// Run and validate.
	if msgs, err := storage.List(); err != nil || len(msgs) != wantRows {
		t.Errorf("error was not expected while listing messages: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMessageListError(t *testing.T) {
	// Mock the database.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// We expect a query to be run.
	mock.ExpectQuery(listSelectQuery).WillReturnError(fmt.Errorf("some error"))

	// Instantiate storage.
	storage := NewMessageStorage(db)

	// Run and validate.
	if _, err := storage.List(); err == nil {
		t.Errorf("error WAS expected while listing messages: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func getMessageRows(count int) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"user_id", "create_date", "message", "is_palindrome", "last_updated_date", "last_updated_by"})
	if count > 0 {
		for i := 0; i < count; i++ {
			itxt := strconv.Itoa(i)
			uid := "userid-" + itxt
			date := "2001-01-02T00:00:00.0Z"
			rows.AddRow(uid, date, "test"+itxt, false, date, uid)
		}
	}
	return rows
}
