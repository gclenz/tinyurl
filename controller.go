package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Controller struct {
	Repository *UrlRepository
}

func NewController(repo *UrlRepository) *Controller {
	return &Controller{
		Repository: repo,
	}
}

func (c *Controller) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func (c *Controller) CreateUrl(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var url UrlData

	err := dec.Decode(&url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = c.Repository.Create(&url, r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"shortId\": \"%s\"}", url.ID)))
}
