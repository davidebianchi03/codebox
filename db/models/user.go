package models

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID                uint           `gorm:"primarykey" json:"-"`
	Email             string         `gorm:"column:email; size:255; unique; not null;" json:"email"`
	Password          string         `gorm:"column:password; not null;" json:"-"`
	FirstName         string         `gorm:"column:first_name; size:255;" json:"first_name"`
	LastName          string         `gorm:"column:last_name; size:255;" json:"last_name"`
	Groups            []Group        `gorm:"many2many:user_groups;" json:"groups"`
	SshPrivateKey     string         `gorm:"column:ssh_private_key; not null;" json:"-"`
	SshPublicKey      string         `gorm:"column:ssh_public_key; not null;" json:"-"`
	IsSuperuser       bool           `gorm:"column:is_superuser; column:is_superuser; default:false" json:"is_superuser"`
	IsTemplateManager bool           `gorm:"column:is_template_manager; default:false" json:"is_template_manager"`
	CreatedAt         time.Time      `json:"-"`
	UpdatedAt         time.Time      `json:"-"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func generateSshKeys() (string, string, error) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "sshkeygen")
	if err != nil {
		return "", "", fmt.Errorf("error creating temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory

	// Define file paths within the temporary directory
	privateKeyPath := filepath.Join(tempDir, "id_rsa")
	publicKeyPath := privateKeyPath + ".pub"

	// Run ssh-keygen command to generate keys
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "2048", "-f", privateKeyPath, "-N", "", "-C", "")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("error generating SSH keys: %w", err)
	}

	// Read the private key file content
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", "", fmt.Errorf("error reading private key file: %w", err)
	}
	privateKey := string(privateKeyBytes)

	// Read the public key file content
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", "", fmt.Errorf("error reading public key file: %w", err)
	}
	publicKey := string(publicKeyBytes)

	return privateKey, publicKey, nil
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if u.SshPrivateKey == "" || u.SshPublicKey == "" {
		u.SshPrivateKey, u.SshPublicKey, err = generateSshKeys()
		if err != nil {
			return fmt.Errorf("failed to create ssh keys: %s", err)
		}
	}

	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

/*
GetLastLogin retrieves the last login time of the user by checking
the most recent token creation time.
If the user has never logged in, it returns nil.
*/
func (u *User) GetLastLogin() (*time.Time, error) {
	var token Token
	result := dbconn.DB.Where("user_id = ?", u.ID).Order("created_at DESC").First(&token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // No login records found
		}
		return nil, result.Error // Some other error occurred
	}
	return &token.CreatedAt, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func ValidatePassword(password string) error {
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasSpecialSymbol := regexp.MustCompile(`[!_\-,.?]`).MatchString(password)

	passwordValid := len(password) >= 10 && hasUppercase && hasSpecialSymbol

	if !passwordValid {
		return errors.New(
			"invalid password, it must be at least 10 characters long and " +
				"include at least one uppercase letter and one special symbol (!_-,.?!)",
		)
	}
	return nil
}

func CreateUser(email, firstName, lastName, password string, isSuperUser, isTemplateManager bool) (user *User, err error) {
	password, err = HashPassword(password)
	if err != nil {
		return nil, err
	}

	// create new user
	newUser := User{
		Email:             email,
		FirstName:         firstName,
		LastName:          lastName,
		Password:          password,
		IsSuperuser:       isSuperUser,
		IsTemplateManager: isTemplateManager,
	}

	r := dbconn.DB.Create(&newUser)
	if r.Error != nil {
		return nil, r.Error
	}

	return user, nil
}

/*
ListUsers retrieves all users from the database ordered by -CreatedAt.
Limit specifies the maximum number of users to retrieve.
If limit is -1 all users are retrieved.
*/
func ListUsers(limit int) (users *[]User, err error) {
	users = &[]User{}
	result := dbconn.DB.Order("created_at DESC").Limit(limit).Find(users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

/*
RetrieveUserByEmail retrieves a user by their email address.
*/
func RetrieveUserByEmail(email string) (user *User, err error) {
	result := dbconn.DB.Where("email=?", email).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected > 0 {
		return user, nil
	}
	return nil, nil
}

/*
CountAllUsers counts the total number of users in the database.
*/
func CountAllUsers() (count int64, err error) {
	if err = dbconn.DB.Model(User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
