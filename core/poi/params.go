package poi

import (
	"github.com/iotaledger/hive.go/app"
)

type ParametersPOI struct {
	BindAddress string `default:"localhost:9687" usage:"bind address on which the POI HTTP server listens"`
}

var ParamsPOI = &ParametersPOI{}

var params = &app.ComponentParams{
	Params: map[string]any{
		"poi": ParamsPOI,
	},
	Masked: nil,
}
