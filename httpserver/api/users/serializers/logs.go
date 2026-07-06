package serializers

import (
	"gitlab.com/codebox4073715/codebox/logging"
)

type SystemLogRowSerializer struct {
	Timestamp string `json:"timestamp"`
	Module    string `json:"module"`
	Function  string `json:"function"`
	Level     string `json:"level"`
	Log       string `json:"log"`
}

func LoadSystemLogRow(l logging.LogRow) SystemLogRowSerializer {
	return SystemLogRowSerializer{
		Timestamp: l.Timestamp,
		Module:    l.Module,
		Function:  l.Function,
		Level:     l.Level,
		Log:       l.Log,
	}
}

func LoadMultipleSystemLogRows(logs []logging.LogRow) []SystemLogRowSerializer {
	serialized := make([]SystemLogRowSerializer, len(logs))
	for i, log := range logs {
		serialized[i] = LoadSystemLogRow(log)
	}
	return serialized
}
