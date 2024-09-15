package db

import "gorm.io/gorm"

const (
	ConnectionTypeHttp = "http"
	ConnectionTypeWS   = "ws"
)

type ForwardedPort struct {
	gorm.Model
	portNumber     uint   `gorm:""`
	active         bool   `gorm:"default:true"`
	connectionType string `gorm:"size:40; not null;default:http;"`
	public         bool   `gorm:"default:false"`
}

func (*ForwardedPort) GetPublicUrl() string {
	return ""
}

func (*ForwardedPort) RequestHeaders() map[string]interface{} {
	var headers map[string]interface{}
	return headers
}
