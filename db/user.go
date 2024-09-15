package db

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id            uint   `gorm:"unique;primaryKey;autoIncrement"`
	Email         string `gorm:"unique;not null;"`
	Password      string `gorm:"not null;"`
	FirstName     string `gorm:""`
	LastName      string `gorm:""`
	SshPrivateKey string `gorm:"not null;"`
	SshPublicKey  string `gorm:"not null;"`
	IsSuperUser   bool   `gorm:"default:false"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func generateSshKeys() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", "", err
	}

	// generate private key
	var privateKeyBuf bytes.Buffer
	privateKeyBufW := io.Writer(&privateKeyBuf)
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	pem.Encode(privateKeyBufW, privateKeyPEM)
	privateKeyStr := privateKeyBuf.String()

	// generate public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubKeyBytes := ssh.MarshalAuthorizedKey(pub)
	publicKeyStr := string(pubKeyBytes)

	return privateKeyStr, publicKeyStr, nil
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	// hash della password se la password è cambiata
	if tx.Statement.Changed("Password") {
		u.Password, err = hashPassword(u.Password)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// hash della password
	u.Password, err = hashPassword(u.Password)
	if err != nil {
		return err
	}
	u.SshPrivateKey, u.SshPublicKey, err = generateSshKeys()
	if err != nil {
		return fmt.Errorf("failed to create ssh keys: %s", err)
	}
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
