package serializers

import "gitlab.com/codebox4073715/codebox/db/models"

type AuthenticationSettingsSerializer struct {
	IsSignUpOpen       bool   `json:"is_signup_open"`
	IsSignUpRestricted bool   `json:"is_signup_restricted"`
	AllowedEmailRegex  string `json:"allowed_emails_regex"`
	BlockedEmailRegex  string `json:"blocked_emails_regex"`
}

func LoadAuthenticationSettingsSerializer(is *models.AuthenticationSettings) *AuthenticationSettingsSerializer {
	if is == nil {
		return nil
	}

	return &AuthenticationSettingsSerializer{
		IsSignUpOpen:       is.IsSignUpOpen,
		IsSignUpRestricted: is.IsSignUpRestricted,
		AllowedEmailRegex:  is.AllowedEmailRegex,
		BlockedEmailRegex:  is.BlockedEmailRegex,
	}
}

type EmailServiceConfiguredSerializer struct {
	IsConfigured bool `json:"is_configured"`
}

func LoadEmailServiceConfiguredSerializer(isConfigured bool) *EmailServiceConfiguredSerializer {
	return &EmailServiceConfiguredSerializer{
		IsConfigured: isConfigured,
	}
}
