package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"gitlab.com/codebox4073715/codebox/api/serializers"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/bgtasks"
	"gitlab.com/codebox4073715/codebox/db/models"
)

// HandleAdminListUsers godoc
// @Summary Admin List Users
// @Schemes
// @Description List all users ordered by creation date descending
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} []serializers.AdminUserSerializer
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

	c.JSON(http.StatusOK, serializers.LoadMultipleAdminUserSerializer(*users))
}

// HandleAdminRetrieveUser godoc
// @Summary Admin Retrieve User
// @Schemes
// @Description Admin Retrieve User
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} serializers.AdminUserSerializer
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

	c.JSON(http.StatusOK, serializers.LoadAdminUserSerializer(user))
}

type AdminCreateUserRequestBody struct {
	Email             string `json:"email" binding:"required,email"`
	Password          string `json:"password" binding:"required"`
	FirstName         string `json:"first_name" binding:"required"`
	LastName          string `json:"last_name" binding:"required"`
	IsSuperuser       bool   `json:"is_superuser"`
	IsTemplateManager bool   `json:"is_template_manager"`
}

// HandleAdminCreateUser godoc
// @Summary Admin Create User
// @Schemes
// @Description Admin Create User
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body AdminCreateUserRequestBody true "User info"
// @Success 200 {object} serializers.AdminUserSerializer
// @Router /api/v1/admin/users [post]
func HandleAdminCreateUser(c *gin.Context) {
	var reqBody AdminCreateUserRequestBody

	if c.ShouldBindBodyWithJSON(&reqBody) != nil {
		utils.ErrorResponse(c, 400, "invalid or missing argument")
		return
	}

	// check if exists another user with the same email address
	user, err := models.RetrieveUserByEmail(reqBody.Email)
	if err != nil {
		utils.ErrorResponse(c, 500, "internal server error")
		return
	}

	if user != nil {
		utils.ErrorResponse(c, 409, "another user with the same email already exists")
		return
	}

	// validate password
	if err := models.ValidatePassword(reqBody.Password); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// create new user
	u, err := models.CreateUser(
		reqBody.Email,
		reqBody.FirstName,
		reqBody.LastName,
		reqBody.Password,
		reqBody.IsSuperuser,
		reqBody.IsTemplateManager,
	)
	if err != nil {
		utils.ErrorResponse(c, 500, "internal server error")
		return
	}

	c.JSON(http.StatusCreated, serializers.LoadAdminUserSerializer(u))
}

type AdminUpdateUserRequestBody struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	IsSuperuser       bool   `json:"is_superuser"`
	IsTemplateManager bool   `json:"is_template_manager"`
}

// HandleAdminUpdateUser godoc
// @Summary Admin update user
// @Schemes
// @Description Admin update user
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body AdminUpdateUserRequestBody true "User info"
// @Success 200 {object} serializers.AdminUserSerializer
// @Router /api/v1/admin/users/{email} [put]
func HandleAdminUpdateUser(c *gin.Context) {
	currentUser, _ := utils.GetUserFromContext(c)
	email, _ := c.Params.Get("email")

	var requestBody AdminUpdateUserRequestBody
	if c.ShouldBindBodyWithJSON(&requestBody) != nil {
		utils.ErrorResponse(c, 400, "invalid or missing argument")
		return
	}

	user, err := models.RetrieveUserByEmail(email)
	if err != nil {
		utils.ErrorResponse(c, 500, "internal server error")
		return
	}

	if user == nil {
		utils.ErrorResponse(c, 404, "user not found")
		return
	}

	// update fields
	user.FirstName = requestBody.FirstName
	user.LastName = requestBody.LastName
	user.IsSuperuser = requestBody.IsSuperuser
	user.IsTemplateManager = requestBody.IsTemplateManager

	// prevent admin from removing their own superuser status
	// this could lock them out of the admin panel
	if !requestBody.IsSuperuser && user.Email == currentUser.Email {
		user.IsSuperuser = true
	}

	if err := models.UpdateUser(user); err != nil {
		utils.ErrorResponse(c, 500, "internal server error")
		return
	}

	c.JSON(http.StatusOK, serializers.LoadAdminUserSerializer(user))
}

// HandleAdminDeleteUser godoc
// @Summary Admin delete user
// @Schemes
// @Description Admin delete user
// @Tags Admin
// @Accept json
// @Produce json
// @Success 204
// @Router /api/v1/admin/users/{email} [delete]
func HandleAdminDeleteUser(c *gin.Context) {
	currentUser, _ := utils.GetUserFromContext(c)
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

	if user == &currentUser {
		utils.ErrorResponse(c, 400, "you cannot delete yourself")
	}

	bgtasks.BgTasksEnqueuer.Enqueue("delete_user", work.Q{"user_email": user.Email})

	c.JSON(http.StatusNoContent, gin.H{
		"detail": "user deletion has been scheduled",
	})
}

type AdminSetUserPasswordRequestBody struct {
	Password string `json:"password" binding:"required"`
}

// AdminSetUserPassword godoc
// @Summary Admin update user password
// @Schemes
// @Description Admin update user password
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body AdminSetUserPasswordRequestBody true "User info"
// @Success 200
// @Router /api/v1/admin/users/{email}/set-password [post]
func HandleAdminSetUserPassword(c *gin.Context) {
	email, _ := c.Params.Get("email")

	var requestBody AdminSetUserPasswordRequestBody
	if c.ShouldBindBodyWithJSON(&requestBody) != nil {
		utils.ErrorResponse(c, 400, "invalid or missing argument")
		return
	}

	user, err := models.RetrieveUserByEmail(email)
	if err != nil {
		utils.ErrorResponse(c, 500, "internal server error")
		return
	}

	if user == nil {
		utils.ErrorResponse(c, 404, "user not found")
		return
	}

	// validate password
	if err := models.ValidatePassword(requestBody.Password); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	user.Password, err = models.HashPassword(requestBody.Password)
	if err != nil {
		utils.ErrorResponse(c, 500, "internal server error")
		return
	}

	if err := models.UpdateUser(user); err != nil {
		utils.ErrorResponse(c, 500, "internal server error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"detail": "password changed",
	})
}

// HandleAdminImpersonateUser godoc
// @Summary API for admins to impersonate a user
// @Schemes
// @Description API for admins to impersonate a user
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200
// @Router /api/v1/admin/users/{email}/impersonate [post]
func HandleAdminImpersonateUser(c *gin.Context) {
	// TODO: limit the impersonation only to jwt tokens used in cookies
	email, _ := c.Params.Get("email")

	user, err := models.RetrieveUserByEmail(email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if user == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	if user.IsSuperuser {
		utils.ErrorResponse(c, http.StatusBadRequest, "cannot impersonate a superuser")
		return
	}

	// start impersonation and create impersonation log
	token, err := utils.GetTokenFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	token.ImpersonatedUserID = user.ID
	token.ImpersonatedUser = user

	if err := models.UpdateToken(token); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	_, err = models.CreateImpersonationLog(
		token,
		token.User,
		c.ClientIP(),
		*user,
	)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"detail": "impersonation started",
	})
}

// HandleStopImpersonation godoc
// @Summary API to stop the impersonation of a user
// @Schemes
// @Description API to stop the impersonation of a user
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200
// @Router /api/v1/admin/users/{email}/stop-impersonation [post]
func HandleStopImpersonation(c *gin.Context) {
	token, err := utils.GetTokenFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if token.ImpersonatedUser == nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"no user is being impersonated in this session",
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"detail": "impersonation has been stopped",
	})
}
