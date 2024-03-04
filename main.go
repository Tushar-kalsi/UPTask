package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"
	"uptask/db"
	"uptask/handlers"
	"uptask/models"
	"uptask/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func CheckAndNotifyTasks(ctx context.Context, db *sql.DB) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done(): // Check if the context is cancelled or expired
			log.Println("Stopping task notification service")
			return // Exit the function
		case <-ticker.C:
			now := time.Now()
			var tasks []models.Task

			// Query tasks due within the next hour (customize the query as needed)
			rows, err := db.Query("SELECT id, title, description, priority, due_date FROM tasks WHERE due_date <= $1", now.Add(1*time.Hour))
			if err != nil {
				log.Println("Error querying tasks:", err)
				continue
			}
			defer rows.Close()

			for rows.Next() {
				var task models.Task
				if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDate); err != nil {
					log.Println("Error scanning task:", err)
					continue
				}
				tasks = append(tasks, task)
			}

			if err := rows.Err(); err != nil {
				log.Println("Error fetching tasks:", err)
				continue
			}

			// Send notifications for due tasks
			for _, task := range tasks {
				utils.SendNotification(task)
			}
		}
	}
}

func main() {

	handlers.DB = db.ConnectDB()
	defer handlers.DB.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure the context is canceled when main exits

	// Start the background service in a goroutine
	go CheckAndNotifyTasks(ctx, handlers.DB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/tasks", handlers.CreateTaskHandler)
	r.Delete("/tasks/{taskID}", handlers.DeleteTaskHandler) // Delete
	r.Get("/tasks", handlers.ListTasksHandler)              // Read all
	r.Put("/tasks/{taskID}", handlers.UpdateTaskHandler)    // Update
	r.Delete("/tasks/{taskID}", handlers.DeleteTaskHandler) // Delete

	http.ListenAndServe(":3000", r)
}
