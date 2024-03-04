package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"uptask/handlers"
	"uptask/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTaskHandler(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set the global DB variable to our mock database
	handlers.DB = db

	// Setup expectations
	mock.ExpectQuery(`INSERT INTO tasks`).
		WithArgs("Test Task", "Test Description", 1, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Create a Task to send in the HTTP request body
	task := models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    1,
		DueDate:     time.Now(),
	}
	taskBytes, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}

	// Create an HTTP request
	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(taskBytes))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateTaskHandler)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	// Check the response body
	var createdTask models.Task
	err = json.Unmarshal(rr.Body.Bytes(), &createdTask)
	assert.NoError(t, err)
	assert.Equal(t, task.Title, createdTask.Title, "Expected the task title to match")

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
