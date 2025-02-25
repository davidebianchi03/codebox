package models

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID            uint    `gorm:"primarykey"`
	Email         string  `gorm:"size:255; unique; not null;"`
	Password      string  `gorm:"not null;"`
	FirstName     string  `gorm:"size:255;"`
	LastName      string  `gorm:"size:255;"`
	Groups        []Group `gorm:"many2many:user_groups;"`
	SshPrivateKey string  `gorm:"not null;"`
	SshPublicKey  string  `gorm:"not null;"`
	IsSuperuser   bool    `gorm:"column:is_superuser; default:false"`
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
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "2048", "-f", privateKeyPath, "-N", "")
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}
