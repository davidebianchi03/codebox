package db

import "gorm.io/gorm"

const (
	ConnectionTypeHttp = "http"
	ConnectionTypeWS   = "ws"
)

type ForwardedPort struct {
	gorm.Model
	portNumber     uint   `gorm:"column:port_number;"`
	active         bool   `gorm:"column:active; default:true"`
	connectionType string `gorm:"column:connection_type; size:40; not null;default:http;"`
	public         bool   `gorm:"column:public; default:false"`
}

func (*ForwardedPort) GetPublicUrl() string {
	return ""
}

func (*ForwardedPort) RequestHeaders() map[string]interface{} {
	var headers map[string]interface{}
	return headers
}
