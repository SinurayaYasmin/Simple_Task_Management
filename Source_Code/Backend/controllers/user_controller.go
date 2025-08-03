package controllers

import (
	"SimpleTaskManager/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateAccount(ctx *gin.Context) {
	var newUser models.User

	err := ctx.ShouldBindJSON(&newUser)

	if err != nil || newUser.Username == "" || newUser.Email == "" || newUser.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing input"})
		return
	}

	SafeUser, err := models.SignUp(&newUser)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Email already used"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusCreated, models.CreateResponse(true, "Created Amccount Success", SafeUser))
}

func SignIn(ctx *gin.Context) {
	var user models.User

	err := ctx.ShouldBindJSON(&user)

	if err != nil || user.Email == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing input"})
		return
	}

	existUser, err := models.SignIn(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or passsword"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(user.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing input"})
		return
	}

	ctx.JSON(http.StatusOK, models.CreateResponse(true, "Sign In Success", existUser.ID))

}

func GetUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil || id == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID Number"})
		return
	}

	existUser, err := models.GetUser(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User Not Found"})
		return
	}

	ctx.JSON(http.StatusOK, models.CreateResponse(true, "Get User Success", existUser))

}

func GetAllUser(ctx *gin.Context) {
	allUser, err := models.GetAllUser()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all users"})
	}

	ctx.JSON(http.StatusOK, models.CreateResponse(true, "Success to get all users", allUser))

}
