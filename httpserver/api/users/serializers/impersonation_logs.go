package serializers

import (
	"time"

	"gitlab.com/codebox4073715/codebox/db/models"
)

type ImpersonationLogSerializer struct {
	ID                      uint                `json:"id"`
	Impersonator            AdminUserSerializer `json:"impersonator"`
	ImpersonatorIPAddress   string              `json:"impersonator_ip_address"`
	ImpersonationStartedAt  time.Time           `json:"impersonation_started_at"`
	ImpersonationFinishedAt *time.Time          `json:"impersonation_finished_at"`
	SessionExpired          bool                `json:"session_expired"`
}

func LoadImpersonationLogSerializer(log *models.ImpersonationLog) *ImpersonationLogSerializer {
	if log == nil {
		return nil
	}

	sessionExpired := false

	if log.Token == nil {
		sessionExpired = true
	} else {
		if time.Since(*log.Token.ExpirationDate) > 0 {
			sessionExpired = true
		}
	}

	return &ImpersonationLogSerializer{
		ID:                      log.ID,
		Impersonator:            *LoadAdminUserSerializer(&log.Impersonator),
		ImpersonatorIPAddress:   log.ImpersonatorIPAddress,
		ImpersonationStartedAt:  log.ImpersonationStartedAt,
		ImpersonationFinishedAt: log.ImpersonationFinishedAt,
		SessionExpired:          sessionExpired,
	}
}

func LoadMultipleImpersonationLogSerializer(logs []models.ImpersonationLog) []ImpersonationLogSerializer {
	serializers := make([]ImpersonationLogSerializer, len(logs))

	for i, l := range logs {
		serializers[i] = *LoadImpersonationLogSerializer(&l)
	}

	return serializers
}
