package proxy

type NginxProxyManager struct {
	NPMEndpoint string
	NPMUser     string
	NPMPassword string
	NPMToken    string
}

func InitNPMInterface(serverEndpoint string, username string, password string) (*NginxProxyManager, error) {
	token, err := NPMLogin(serverEndpoint, username, password)
	if err != nil {
		return nil, err
	}
	proxyConfig := NginxProxyManager{
		NPMEndpoint: serverEndpoint,
		NPMUser:     username,
		NPMPassword: password,
		NPMToken:    token,
	}
	return &proxyConfig, nil
}
