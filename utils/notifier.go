package utils

import (
	"fmt"
	"uptask/models"
)

func SendNotification(task models.Task) {
	fmt.Printf("Sending notification: Task '%s' is due at %s.\n", task.Title, task.DueDate)
}
