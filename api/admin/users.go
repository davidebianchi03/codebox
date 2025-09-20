package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/serializers"
	"gitlab.com/codebox4073715/codebox/api/utils"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
)

// HandleAdminListUsers godoc
// @Summary Admin List Users
// @Schemes
// @Description List all users ordered by creation date descending
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} []serializers.UserSerializer
// @Router /api/v1/admin/users [get]
func HandleAdminListUsers(c *gin.Context) {
	limit := c.Query("limit")
	if limit == "" {
		limit = "-1"
	}

	// validate limit
	parsedLimit, err := strconv.Atoi(limit)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid limit")
		return
	}

	if parsedLimit < -1 || parsedLimit == 0 {
		utils.ErrorResponse(c, 400, "invalid limit")
		return
	}

	// retrieve users
	users, err := models.ListUsers(parsedLimit)
	if err != nil {
		utils.ErrorResponse(c, 500, "internal server error")
		return
	}

	c.JSON(http.StatusOK, serializers.LoadMultipleUserSerializer(*users))
}

// HandleAdminRetrieveUser godoc
// @Summary Admin Retrieve User
// @Schemes
// @Description Admin Retrieve User
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} serializers.UserSerializer
// @Router /api/v1/admin/users/{email} [get]
func HandleAdminRetrieveUser(c *gin.Context) {
	email, _ := c.Params.Get("email")

	user, err := models.RetrieveUserByEmail(email)
	if err != nil {
		utils.ErrorResponse(c, 500, "internal server error")
		return
	}

	if user == nil {
		utils.ErrorResponse(c, 404, "user not found")
		return
	}

	c.JSON(http.StatusOK, serializers.LoadUserSerializer(user))
}

func HandleAdminCreateUser(c *gin.Context) {
	var reqBody struct {
		Email             string `json:"email" binding:"required,email"`
		Password          string `json:"password" binding:"required"`
		FirstName         string `json:"first_name" binding:"required"`
		LastName          string `json:"last_name" binding:"required"`
		IsSuperuser       bool   `json:"is_superuser"`
		IsTemplateManager bool   `json:"is_template_manager"`
	}

	if c.ShouldBindBodyWithJSON(&reqBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing or invalid field",
		})
		return
	}

	// check if exists another user with the same email address
	users := []models.User{}
	r := dbconn.DB.Find(&users, map[string]interface{}{
		"email": reqBody.Email,
	})

	if r.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if len(users) > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"detail": "another user with the same email already exists",
		})
		return
	}

	// validate password
	if err := models.ValidatePassword(reqBody.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	password, err := models.HashPassword(reqBody.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	// create new user
	newUser := models.User{
		Email:             reqBody.Email,
		FirstName:         reqBody.FirstName,
		LastName:          reqBody.LastName,
		Password:          password,
		IsSuperuser:       reqBody.IsSuperuser,
		IsTemplateManager: reqBody.IsTemplateManager,
	}

	r = dbconn.DB.Create(&newUser)
	if r.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func HandleAdminUpdateUser(c *gin.Context) {
	currentUser, _ := utils.GetUserFromContext(c)

	var user *models.User
	email, _ := c.Params.Get("email")

	if dbconn.DB.Find(&user, map[string]interface{}{
		"email": email,
	}).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "user not found",
		})
		return
	}

	var requestBody struct {
		FirstName         *string `json:"first_name"`
		LastName          *string `json:"last_name"`
		IsSuperuser       *bool   `json:"is_superuser"`
		IsTemplateManager *bool   `json:"is_template_manager"`
	}

	if c.ShouldBindBodyWithJSON(&requestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "invalid or missing argument",
		})
		return
	}

	if requestBody.FirstName != nil {
		user.FirstName = *requestBody.FirstName
	}

	if requestBody.LastName != nil {
		user.LastName = *requestBody.LastName
	}

	if requestBody.IsSuperuser != nil && user.Email != currentUser.Email {
		user.IsSuperuser = *requestBody.IsSuperuser
	}

	if requestBody.IsTemplateManager != nil {
		user.IsTemplateManager = *requestBody.IsTemplateManager
	}

	dbconn.DB.Save(&user)

	c.JSON(http.StatusOK, user)
}

func HandleAdminSetUserPassword(c *gin.Context) {
	var user *models.User
	var err error
	email, _ := c.Params.Get("email")

	if dbconn.DB.Find(&user, map[string]interface{}{
		"email": email,
	}).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "user not found",
		})
		return
	}

	var requestBody struct {
		Password string `json:"password" binding:"required"`
	}

	if c.ShouldBindBodyWithJSON(&requestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "invalid or missing argument",
		})
		return
	}

	// validate password
	if err := models.ValidatePassword(requestBody.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	user.Password, err = models.HashPassword(requestBody.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	dbconn.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"detail": "password changed",
	})
}
