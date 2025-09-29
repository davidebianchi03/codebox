package models

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gorm.io/gorm"
)

type Runner struct {
	ID                 uint           `gorm:"primarykey" json:"id"`
	Name               string         `gorm:"column:name; size:255;unique;not null;" json:"name"`
	Token              string         `gorm:"column:token; size:255;unique;not null;" json:"-"`
	Port               uint           `gorm:"column:port; default:0;" json:"-"`
	Type               string         `gorm:"column:type; size:255;" json:"type"`
	Restricted         bool           `gorm:"column:restricted; default:false;" json:"-"`
	AllowedGroups      []Group        `gorm:"many2many:runner_allowed_groups;" json:"-"`
	UsePublicUrl       bool           `gorm:"column:use_public_url; default:false;" json:"use_public_url"`
	PublicUrl          string         `gorm:"column:public_url; type:text;" json:"public_url"`
	LastContact        *time.Time     `gorm:"column:last_contact;" json:"last_contact"`
	Version            string         `gorm:"column:version; default:''; size:255;" json:"version"`
	DeletionInProgress bool           `gorm:"column:deletion_in_progress;default:false;not null;"`
	CreatedAt          time.Time      `gorm:"column:created_at;" json:"-"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;" json:"-"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

/*
ListRunners retrieves a list of runners with pagination support.
If limit is -1, it retrieves all runners.
*/
func ListRunners(limit int, offset int) ([]Runner, error) {
	var runners []Runner
	if err := dbconn.DB.Limit(limit).Offset(offset).Find(&runners).Error; err != nil {
		return nil, err
	}
	return runners, nil
}

/*
RetrieveRunnerByID retrieves a runner by its ID
*/
func RetrieveRunnerByID(id uint) (*Runner, error) {
	var runner Runner
	if err := dbconn.DB.
		Preload("AllowedGroups").
		First(&runner, map[string]interface{}{
			"ID": id,
		}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &runner, nil
}

/*
RetrieveRunnerByName retrieves a runner by its name
*/
func RetrieveRunnerByName(name string) (*Runner, error) {
	var runner Runner
	if err := dbconn.DB.
		Preload("AllowedGroups").
		First(&runner, map[string]interface{}{
			"Name": name,
		}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return &runner, nil
}

/*
DoesRunnerExistWithUrl checks if a runner with the given public url exists
*/
func DoesRunnerExistWithUrl(url string) (bool, error) {
	exists := false
	err := dbconn.DB.Model(Runner{}).
		Select("count(*) > 0").
		Where("public_url = ?", url).
		Find(&exists).
		Error
	return exists, err
}

/*
CountOnlineRunners counts the number of online runners.
A runner is considered online if its last contact time is within the last 5 minutes.
*/
func CountOnlineRunners() (int64, error) {
	var count int64
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
	if err := dbconn.DB.Model(&Runner{}).
		Where("last_contact >= ?", fiveMinutesAgo).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

/*
Generate a random string long as the input param
*/
func generateToken(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-#!_=+"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

/*
Create a new runner
*/
func CreateRunner(
	runnerName string,
	runnerType string,
	usePublicUrl bool,
	publicUrl string,
) (*Runner, error) {
	// generate the token
	token := ""
	exists := false
	for ok := true; ok; ok = exists {
		token = fmt.Sprintf("cbrt-%s", generateToken(30))

		if err := dbconn.DB.Model(Runner{}).
			Select("count(*) > 0").
			Where("token = ?", token).
			Find(&exists).
			Error; err != nil {
			return nil, err
		}
	}

	fmt.Println("token: " + token)

	runner := Runner{
		Name:         runnerName,
		Type:         runnerType,
		Token:        token,
		UsePublicUrl: usePublicUrl,
		PublicUrl:    publicUrl,
	}

	if err := dbconn.DB.Create(&runner).Error; err != nil {
		return nil, err
	}

	return &runner, nil
}
