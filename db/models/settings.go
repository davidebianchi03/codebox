package models

import "time"

type AuthenticationSettings struct {
	SingletonModel
	IsSignUpOpen                bool   `gorm:"column:is_signup_open; default:false"`
	IsSignUpRestricted          bool   `gorm:"column:is_signup_restricted; default:false"`
	AllowedEmailRegex           string `gorm:"column:allowed_email_regex; type:text;"`
	BlockedEmailRegex           string `gorm:"column:blocked_email_regex; type:text;"`
	UsersMustBeApproved         bool   `gorm:"column:users_must_be_approved; default:false"`
	ApprovedByDefaultEmailRegex string `gorm:"column:approved_by_default_email_regex; type:text;"`
}

type AnalyticsConfig struct {
	SingletonModel
	SendAnalyticsData      bool       `gorm:"column:send_analytics_data; default:false"`
	AnalyticsBannerSent    bool       `gorm:"column:analytics_banner_sent; default:false"`
	LastAttempt            *time.Time `gorm:"column:last_attempt"`
	LastSuccessfullAttempt *time.Time `gorm:"column:last_successfull_attempt"`
}
