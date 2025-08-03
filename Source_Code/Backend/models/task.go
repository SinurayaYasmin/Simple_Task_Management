package models

import (
	"SimpleTaskManager/config"
	"fmt"

	"github.com/google/uuid"

	"database/sql"
	"time"
)

type NewTask struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Deadline    time.Time `json:"deadline"`
	Created     time.Time `json:"created"`
	CreatedBy   string    `json:"created_by"`
}

type DetailTask struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Deadline    time.Time  `json:"deadline"`
	Created     *time.Time `json:"created"`
	Updated     *time.Time `json:"updated"`
	Finished    *time.Time `json:"finished"`
	Deleted     *time.Time `json:"deleted"`
	CreatedBy   *string    `json:"created_by"`
}

type Assignee struct {
	TaskID    uuid.UUID `json:"task_id"`
	TaskTitle string    `json: "task_title`
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	Username  string    `json: "username`
}

type AssigneeInput struct {
	TaskID    uuid.UUID   `json:"task_id" binding:"required"`
	TaskTitle string      `json: "task_title`
	UserID    []uuid.UUID `json:"user_id" binding:"required"`
	Username  string      `json: "username`
}

// Create Task
func CreateTask(task *NewTask) (*NewTask, error) {
	var newTask NewTask
	query := `INSERT INTO tasks (title, description, deadline, created_at, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id, title, description, status, deadline, created_at, created_by`
	row := config.DB.QueryRow(query, task.Title, task.Description, task.Deadline, task.Created, task.CreatedBy)
	err := row.Scan(&newTask.ID, &newTask.Title, &newTask.Description, &newTask.Status, &newTask.Deadline, &newTask.Created, &newTask.CreatedBy)

	if err != nil {
		return nil, err
	}

	return &newTask, nil
}

// Get Task
func GetTask(id uuid.UUID) (*DetailTask, error) {
	var task DetailTask
	query := `SELECT * FROM tasks WHERE id = $1`
	row := config.DB.QueryRow(query, id)
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Deadline, &task.Created, &task.Updated, &task.Deleted, &task.Finished, &task.CreatedBy)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}

	return &task, nil

}

// Get All Task
func GetAllTask() ([]DetailTask, error) {
	var allTask []DetailTask
	query := `SELECT * FROM tasks`
	row, err := config.DB.Query(query)

	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		var task DetailTask
		err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Deadline, &task.Created, &task.Updated, &task.Deleted, &task.Finished, &task.CreatedBy)

		if err != nil {
			return nil, err
		}

		allTask = append(allTask, task)
	}

	return allTask, nil
}

// Update Task
func UpdateTask(task *DetailTask) (*DetailTask, error) {
	var updatedTask DetailTask
	query := `UPDATE tasks SET title = $1, description = $2, status = $3, deadline = $4, updated_at = $5 WHERE id = $6 RETURNING *`
	row := config.DB.QueryRow(query, task.Title, task.Description, task.Status, task.Deadline, task.Updated, task.ID)

	err := row.Scan(&updatedTask.ID, &updatedTask.Title, &updatedTask.Description, &updatedTask.Status, &updatedTask.Deadline, &updatedTask.Created, &updatedTask.Updated, &updatedTask.Deleted, &updatedTask.Finished, &updatedTask.CreatedBy)

	if err != nil {
		return nil, err
	}

	return &updatedTask, err
}

// Delete Task
func DeleteTask(id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := config.DB.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func FinishedTask(task DetailTask) (*DetailTask, error) {
	var finishedTask DetailTask
	query := `UPDATE tasks set status = 'Done', finished_at = $1 WHERE id = $2 RETURNING *`
	row := config.DB.QueryRow(query, task.Finished, task.ID)

	err := row.Scan(&finishedTask.ID, &finishedTask.Title, &finishedTask.Description, &finishedTask.Status, &finishedTask.Deadline, &finishedTask.Created, &finishedTask.Updated, &finishedTask.Deleted, &finishedTask.Finished, &finishedTask.CreatedBy)

	if err != nil {
		return nil, err
	}

	return &finishedTask, nil
}

// Delete Task
// Choose Assignee

func ChooseAssignee(assignee AssigneeInput) ([]Assignee, error) {
	var choosen []Assignee
	query := `INSERT INTO task_assignee values ($1, $2) RETURNING *`

	for _, userID := range assignee.UserID {
		_, err := config.DB.Exec(query, assignee.TaskID, userID)
		if err != nil {
			return nil, err
		}

	}

	joinQuery := `SELECT ta.task_id, t.title, ta.user_id, u.username FROM task_assignee ta
	JOIN tasks t ON ta.task_id = t.id
	JOIN users u ON ta.user_id = u.id
	WHERE ta.task_id = $1`

	rows, err := config.DB.Query(joinQuery, assignee.TaskID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a Assignee
		err := rows.Scan(&a.TaskID, &a.TaskTitle, &a.UserID, &a.Username)
		if err != nil {
			return nil, err
		}
		choosen = append(choosen, a)
	}

	return choosen, nil
}
