package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/models"
	"github.com/gophish/gophish/webhook"
	"github.com/gorilla/mux"
)

// Webhooks returns a list of webhooks, both active and disabled
func (as *Server) Webhooks(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		whs, err := models.GetWebhooks()
		if err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, whs, http.StatusOK)

	case r.Method == "POST":
		wh := models.Webhook{}
		err := json.NewDecoder(r.Body).Decode(&wh)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid JSON structure"}, http.StatusBadRequest)
			return
		}
		err = models.PostWebhook(&wh)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		JSONResponse(w, wh, http.StatusCreated)
	}
}

// Webhook returns details of a single webhook specified by "id" parameter
func (as *Server) Webhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	wh, err := models.GetWebhook(id)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Webhook not found"}, http.StatusNotFound)
		return
	}
	switch {
	case r.Method == "GET":
		JSONResponse(w, wh, http.StatusOK)

	case r.Method == "DELETE":
		err = models.DeleteWebhook(id)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		log.Infof("Deleted webhook with id: %d", id)
		JSONResponse(w, models.Response{Success: true, Message: "Webhook deleted Successfully!"}, http.StatusOK)

	case r.Method == "PUT":
		wh2 := models.Webhook{}
		err = json.NewDecoder(r.Body).Decode(&wh2)
		wh2.Id = id
		err = models.PutWebhook(&wh2)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		JSONResponse(w, wh2, http.StatusOK)
	}
}

// ValidateWebhook makes an HTTP request to a specified remote url to ensure that it's valid.
func (as *Server) ValidateWebhook(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "POST":
		vars := mux.Vars(r)
		id, _ := strconv.ParseInt(vars["id"], 0, 64)
		wh, err := models.GetWebhook(id)
		if err != nil {
			log.Error(err)
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		err = webhook.Send(webhook.EndPoint{URL: wh.URL, Secret: wh.Secret}, "")
		if err == nil {
			JSONResponse(w, wh, http.StatusOK)
		} else {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
		}
	}
}
