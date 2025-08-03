package controllers

import (
	"SimpleTaskManager/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Create Task
func CreateTask(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil || id == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID Number"})
		return
	}

	var newTask models.NewTask
	err = ctx.ShouldBindJSON(&newTask)

	if err != nil || newTask.Title == "" || newTask.Description == "" || newTask.Status == "" || newTask.Deadline.IsZero() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing input"})
		return
	}

	user, err := models.GetUser(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User Not Found"})
		return
	}

	newTask.CreatedBy = user.Username
	newTask.Created = time.Now()
	task, err := models.CreateTask(&newTask)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusCreated, models.CreateResponse(true, "Created New Task Success", task))

}

// Get Task
func GetTask(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil || id == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID Number"})
		return
	}

	task, err := models.GetTask(id)
	if err != nil {
		if err.Error() == "task not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task Not Found"})
		} else {
			fmt.Println("GetTask error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, models.CreateResponse(true, "Get Task Detail Success", task))

}

// Get All Task
func GetAllTask(ctx *gin.Context) {
	allTask, err := models.GetAllTask()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all task"})
	}

	ctx.JSON(http.StatusOK, models.CreateResponse(true, "Success to get all task", allTask))

}

// Update Task
func UpdateTask(ctx *gin.Context) {
	var task models.DetailTask

	err := ctx.ShouldBindJSON(&task)

	if err != nil || task.Status == "" || task.Title == "" || task.Description == "" || task.Deadline.IsZero() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing input"})
		return
	}
	now := time.Now()
	task.Updated = &now
	updatedTask, err := models.UpdateTask(&task)
	if err != nil {
		fmt.Println("UpdateTask error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, models.CreateResponse(true, "Sign In Success", updatedTask))
}

func FinishedTask(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil || id == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID Number"})
		return
	}

	task, err := models.GetTask(id)
	if err != nil {
		if err.Error() == "task not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task Not Found"})
		} else {
			fmt.Println("GetTask error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}
	now := time.Now()
	task.Finished = &now
	finishedTask, err := models.FinishedTask(*task)

	ctx.JSON(http.StatusOK, models.CreateResponse(true, "Finished Task Success", finishedTask))
}

func Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil || id == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID Number"})
		return
	}

	err = models.DeleteTask(id)
	if err != nil {
		fmt.Println("DeleteTask error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	ctx.JSON(http.StatusOK, models.CreateResponse(true, "Delete Task Success", "-"))
}

// Choose Assignee
func ChooseAssignee(ctx *gin.Context) {
	var assignee models.AssigneeInput

	err := ctx.ShouldBindJSON(&assignee)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}

	if assignee.TaskID == uuid.Nil || len(assignee.UserID) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task ID or User IDs cannot be empty"})
		return
	}

	assigned, err := models.ChooseAssignee(assignee)

	if err != nil {
		fmt.Println("Assignee error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, models.CreateResponse(true, "Assignee Success", assigned))
}
