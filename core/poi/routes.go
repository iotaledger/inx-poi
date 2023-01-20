package poi

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/iotaledger/inx-app/pkg/httpserver"
)

const (
	APIRoute = "poi/v1"

	// ParameterBlockID is used to identify a block by its ID.
	ParameterBlockID = "blockID"

	RouteCreateProof   = "/create/:" + ParameterBlockID
	RouteValidateProof = "/validate"
)

func setupRoutes(e *echo.Echo) {

	e.GET(RouteCreateProof, func(c echo.Context) error {
		resp, err := createProof(c)
		if err != nil {
			return err
		}

		return httpserver.JSONResponse(c, http.StatusOK, resp)
	})

	e.POST(RouteValidateProof, func(c echo.Context) error {
		resp, err := validateProof(c)
		if err != nil {
			return err
		}

		return httpserver.JSONResponse(c, http.StatusOK, resp)
	})
}
