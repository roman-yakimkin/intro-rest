package handlers

import (
	"intro-rest/internal/app/interfaces"
	"net/http"
)

type MemDataController struct {
	store interfaces.Store
}

func NewMemDataController(store interfaces.Store) *MemDataController {
	return &MemDataController{
		store: store,
	}
}

func (c *MemDataController) GetAll(w http.ResponseWriter, r *http.Request) {
	entities, err := c.store.MemData().GetAll()
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	returnSuccessResponse(w, r, "", entities)
}
