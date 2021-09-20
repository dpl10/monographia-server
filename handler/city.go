package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dpl10/monographia-server/model"
)

type (
	cityOutputArray struct {
		Cities []model.City `json:"cities"`
	}
)

// GetCity handler function
func (h *Handler) GetCity(c echo.Context) (err error) {
	x := &cityOutputArray{
		Cities: h.Model.SelectCities(),
	}
	return c.JSON(http.StatusOK, x)
}
