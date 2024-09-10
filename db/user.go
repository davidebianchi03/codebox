package db

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id            uint   `gorm:"unique;primaryKey;autoIncrement"`
	Email         string `gorm:"unique;not null;"`
	Password      string `gorm:"not null;"`
	FirstName     string
	LastName      string
	SshPrivateKey string `gorm:"not null;"`
	SshPublicKey  string `gorm:"not null;"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
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
	// creo la coppia ssh public/private key
	u.SshPrivateKey = "SshPrivateKey"
	u.SshPublicKey = "SshPublicKey"
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
