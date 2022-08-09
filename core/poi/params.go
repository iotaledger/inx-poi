package poi

import (
	"github.com/iotaledger/hive.go/core/app"
)

type ParametersPOI struct {
	BindAddress string `default:"localhost:9687" usage:"bind address on which the POI HTTP server listens"`

	// DebugRequestLoggerEnabled defines whether the debug logging for requests should be enabled
	DebugRequestLoggerEnabled bool `default:"false" usage:"whether the debug logging for requests should be enabled"`
}

var ParamsPOI = &ParametersPOI{}

var params = &app.ComponentParams{
	Params: map[string]any{
		"poi": ParamsPOI,
	},
	Masked: nil,
}
