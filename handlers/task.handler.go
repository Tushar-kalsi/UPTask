package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"uptask/models"

	"github.com/go-chi/chi"
)

var DB *sql.DB

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t models.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO tasks (title, description, priority, due_date) VALUES ($1, $2, $3, $4) RETURNING id`
	err := DB.QueryRow(query, t.Title, t.Description, t.Priority, t.DueDate).Scan(&t.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")

	query := `DELETE FROM tasks WHERE id = $1`
	_, err := DB.Exec(query, taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")

	var task models.Task
	query := `SELECT id, title, description, priority, due_date FROM tasks WHERE id = $1`
	err := DB.QueryRow(query, taskID).Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDate)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	var t models.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE tasks SET title = $1, description = $2, priority = $3, due_date = $4 WHERE id = $5`
	_, err := DB.Exec(query, t.Title, t.Description, t.Priority, t.DueDate, taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	query := `SELECT id, title, description, priority, due_date FROM tasks`
	rows, err := DB.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
