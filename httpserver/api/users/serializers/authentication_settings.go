package serializers

import (
	"encoding/json"

	"gitlab.com/codebox4073715/codebox/db/models"
)

type AuthenticationSettingsSerializer struct {
	IsSignUpOpen                bool   `json:"is_signup_open"`
	IsSignUpRestricted          bool   `json:"is_signup_restricted"`
	AllowedEmailRegex           string `json:"allowed_emails_regex"`
	BlockedEmailRegex           string `json:"blocked_emails_regex"`
	UsersMustBeApproved         bool   `json:"users_must_be_approved"`
	ApprovedByDefaultEmailRegex string `json:"approved_by_default_emails_regex"`
}

func LoadAuthenticationSettingsSerializer(s *models.AuthenticationSettings) *AuthenticationSettingsSerializer {
	if s == nil {
		return nil
	}

	return &AuthenticationSettingsSerializer{
		IsSignUpOpen:                s.IsSignUpOpen,
		IsSignUpRestricted:          s.IsSignUpRestricted,
		AllowedEmailRegex:           s.AllowedEmailRegex,
		BlockedEmailRegex:           s.BlockedEmailRegex,
		UsersMustBeApproved:         s.UsersMustBeApproved,
		ApprovedByDefaultEmailRegex: s.ApprovedByDefaultEmailRegex,
	}
}

func AuthenticationSettingsSerializerFromJSON(data string) (AuthenticationSettingsSerializer, error) {
	var as AuthenticationSettingsSerializer
	if err := json.Unmarshal([]byte(data), &as); err != nil {
		return AuthenticationSettingsSerializer{}, err
	}
	return as, nil
}

type EmailServiceConfiguredSerializer struct {
	IsConfigured bool `json:"is_configured"`
}

func LoadEmailServiceConfiguredSerializer(isConfigured bool) *EmailServiceConfiguredSerializer {
	return &EmailServiceConfiguredSerializer{
		IsConfigured: isConfigured,
	}
}
