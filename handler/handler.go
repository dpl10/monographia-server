package handler

import (
	"github.com/dpl10/monographia-server/model"
)

type (
	// Handler makes interface functions available to all
	Handler struct {
		Model modelInterface
	}
	modelInterface interface {
		SelectCities() []model.City
	}
)

// NewHandler combines all interface functions
func NewHandler(x modelInterface) *Handler {
	return &Handler{
		Model: x,
	}
}
