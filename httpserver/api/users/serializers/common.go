package serializers

import "gitlab.com/codebox4073715/codebox/config"

type VersionSerializer struct {
	Version string `json:"version"`
}

func GetVersionSerializedResponse() VersionSerializer {
	return VersionSerializer{
		Version: config.ServerVersion,
	}
}
