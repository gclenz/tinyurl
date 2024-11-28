package main

import "net/http"

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

	w.WriteHeader(http.StatusCreated)
	w.Write(nil)
}
