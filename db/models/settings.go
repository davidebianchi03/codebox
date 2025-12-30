package models

type AuthenticationSettings struct {
	*SingletonModel
	IsSignUpOpen       bool   `gorm:"column:is_signup_open; default:false"`
	IsSignUpRestricted bool   `gorm:"column:is_signup_restricted; default:false"`
	AllowedEmailRegex  string `gorm:"column:allowed_email_regex; type:text;"`
	BlockedEmailRegex  string `gorm:"column:blocked_email_regex; type:text;"`
}
