package serializers

import (
	"time"

	"gitlab.com/codebox4073715/codebox/db/models"
)

// current user
type CurrentUserSerializer struct {
	Email             string  `json:"email"`
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	IsSuperUser       bool    `json:"is_superuser"`
	IsTemplateManager bool    `json:"is_template_manager"`
	LastLogin         *string `json:"last_login"`
	CreatedAt         string  `json:"created_at"`
	Impersonated      bool    `json:"impersonated"`
}

func LoadCurrentUserSerializer(user *models.User, impersonated bool) *CurrentUserSerializer {
	var lastLoginPtr *string
	if user.ID > 0 {
		lastLogin, err := user.GetLastLogin()
		if err != nil {
			lastLogin = nil
		}

		if lastLogin != nil {
			isoString := lastLogin.Format(time.RFC3339)
			lastLoginPtr = &isoString
		}
	}

	return &CurrentUserSerializer{
		Email:             user.Email,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		IsSuperUser:       user.IsSuperuser,
		IsTemplateManager: user.IsTemplateManager,
		LastLogin:         lastLoginPtr,
		CreatedAt:         user.CreatedAt.Format(time.RFC3339),
		Impersonated:      impersonated,
	}
}

// common
type UserSerializer struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	LastLogin *string `json:"last_login"`
}

func LoadUserSerializer(user *models.User) *UserSerializer {
	var lastLoginPtr *string
	if user.ID > 0 {
		lastLogin, err := user.GetLastLogin()
		if err != nil {
			lastLogin = nil
		}

		if lastLogin != nil {
			isoString := lastLogin.Format(time.RFC3339)
			lastLoginPtr = &isoString
		}
	}

	return &UserSerializer{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		LastLogin: lastLoginPtr,
	}
}

func LoadMultipleUserSerializer(users []models.User) []UserSerializer {
	serializers := make([]UserSerializer, len(users))
	for i, user := range users {
		serializers[i] = *LoadUserSerializer(&user)
	}
	return serializers
}

// admin
type AdminUserSerializer struct {
	Email              string  `json:"email"`
	FirstName          string  `json:"first_name"`
	LastName           string  `json:"last_name"`
	IsSuperUser        bool    `json:"is_superuser"`
	IsTemplateManager  bool    `json:"is_template_manager"`
	DeletionInProgress bool    `json:"deletion_in_progress"`
	LastLogin          *string `json:"last_login"`
	CreatedAt          string  `json:"created_at"`
}

func LoadAdminUserSerializer(user *models.User) *AdminUserSerializer {
	var lastLoginPtr *string
	if user.ID > 0 {
		lastLogin, err := user.GetLastLogin()
		if err != nil {
			lastLogin = nil
		}

		if lastLogin != nil {
			isoString := lastLogin.Format(time.RFC3339)
			lastLoginPtr = &isoString
		}
	}

	return &AdminUserSerializer{
		Email:              user.Email,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		IsSuperUser:        user.IsSuperuser,
		IsTemplateManager:  user.IsTemplateManager,
		LastLogin:          lastLoginPtr,
		CreatedAt:          user.CreatedAt.Format(time.RFC3339),
		DeletionInProgress: user.DeletionInProgress,
	}
}

func LoadMultipleAdminUserSerializer(users []models.User) []AdminUserSerializer {
	serializers := make([]AdminUserSerializer, len(users))
	for i, user := range users {
		serializers[i] = *LoadAdminUserSerializer(&user)
	}
	return serializers
}
