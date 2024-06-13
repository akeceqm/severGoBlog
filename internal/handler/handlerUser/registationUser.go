package handlerUser

import (
	"net/http"
	"post/internal/middlewares"
	"post/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func PostHandleRegistrationUser(c *gin.Context, db *sqlx.DB) {

	var emailRequest struct {
		Email string `json:"email"`
		Name  string `json:"nick_name"`
	}

	errors := make(map[string]string)

	if err := c.BindJSON(&emailRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "указаны не все параметры"})
		return
	}

	// Проверяем, что имя не пустое
	if emailRequest.Name == "" {
		errors["nick_name"] = "Имя пользователя не может быть пустым."
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty name"})
		return
	}

	// Проверяем валидность email
	if !middlewares.ValidateEmail(emailRequest.Email) {
		errors["email"] = "Адрес электронной почты не может быть пустым."
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}

	user, err := services.PostUser(db, emailRequest.Email, emailRequest.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"login": user.Login, "password": user.PasswordHash})
}
