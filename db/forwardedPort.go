package db

import "gorm.io/gorm"

const (
	ConnectionTypeHttp = "http"
	ConnectionTypeWS   = "ws"
)

type ForwardedPort struct {
	gorm.Model
	PortNumber     uint   `gorm:"column:port_number;"`
	Active         bool   `gorm:"column:active; default:true"`
	ConnectionType string `gorm:"column:connection_type; size:40; not null;default:http;"`
	Public         bool   `gorm:"column:public; default:false"`
}

func (*ForwardedPort) GetPublicUrl() string {
	return ""
}

func (*ForwardedPort) RequestHeaders() map[string]interface{} {
	var headers map[string]interface{}
	return headers
}
