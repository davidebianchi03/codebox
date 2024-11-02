package db

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
	Email         string `gorm:"column:email; unique;not null;"`
	Password      string `gorm:"column:password; not null;"`
	FirstName     string `gorm:"column:first_name;"`
	LastName      string `gorm:"column:last_name;"`
	SshPrivateKey string `gorm:"column:ssh_private_key; not null;"`
	SshPublicKey  string `gorm:"column:ssh_public_key; not null;"`
	IsSuperuser   bool   `gorm:"column:is_superuser; default:false"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
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
	// hash della password se la password è cambiata
	if tx.Statement.Changed("Password") {
		u.Password, err = hashPassword(u.Password)
		if err != nil {
			return err
		}
	}

	if u.SshPrivateKey == "" || u.SshPublicKey == "" {
		u.SshPrivateKey, u.SshPublicKey, err = generateSshKeys()
		if err != nil {
			return fmt.Errorf("failed to create ssh keys: %s", err)
		}
	}

	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Password, err = hashPassword(u.Password)
	if err != nil {
		return err
	}
	return u.BeforeSave(tx)
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
