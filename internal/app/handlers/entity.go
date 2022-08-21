package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"intro-rest/internal/app/errors"
	"intro-rest/internal/app/interfaces"
	"intro-rest/internal/app/models"
	"io"
	"net/http"
)

type EntityController struct {
	store interfaces.Store
}

func NewEntityController(store interfaces.Store) *EntityController {
	return &EntityController{
		store: store,
	}
}

func (c *EntityController) Update(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	updEntity := new(models.Entity)
	err = json.Unmarshal(body, updEntity)
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	entity, err := c.store.Entity().GetByID(updEntity.ID)
	if returnErrorResponse(err != nil, w, r, http.StatusNotFound, err, "") {
		return
	}
	entity.Name = updEntity.Name
	_, err = c.store.Entity().Save(entity)
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	returnSuccessResponse(w, r, "entity has been updated", struct {
		Id string `json:"id"`
	}{Id: entity.ID})
}

func (c *EntityController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entityId, ok := vars["id"]
	if returnErrorResponse(!ok, w, r, http.StatusNotFound, errors.ErrEntityNotFound, "") {
		return
	}
	_, err := c.store.Entity().GetByID(entityId)
	if returnErrorResponse(err != nil, w, r, http.StatusNoContent, err, "") {
		return
	}
	err = c.store.Entity().Delete(entityId)
	if returnErrorResponse(err != nil, w, r, http.StatusNoContent, err, "") {
		return
	}
	returnSuccessResponse(w, r, "entity has been deleted", struct {
		Id string `json:"id"`
	}{Id: entityId})
}

func (c *EntityController) GetAll(w http.ResponseWriter, r *http.Request) {
	entities, err := c.store.Entity().GetAll()
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	returnSuccessResponse(w, r, "", entities)
}
