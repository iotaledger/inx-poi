package poi

import (
	"net/http"

	"github.com/labstack/echo/v4"

	restapipkg "github.com/iotaledger/hornet/v2/pkg/restapi"
)

const (
	APIRoute = "poi/v1"

	RouteCreateProof   = "/create/:" + restapipkg.ParameterBlockID
	RouteValidateProof = "/validate"
)

func setupRoutes(e *echo.Echo) {

	e.GET(RouteCreateProof, func(c echo.Context) error {
		resp, err := createProof(c)
		if err != nil {
			return err
		}

		return restapipkg.JSONResponse(c, http.StatusOK, resp)
	})

	e.POST(RouteValidateProof, func(c echo.Context) error {
		resp, err := validateProof(c)
		if err != nil {
			return err
		}
		return restapipkg.JSONResponse(c, http.StatusOK, resp)
	})
}
