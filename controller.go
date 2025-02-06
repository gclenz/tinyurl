package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"
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
	w.Write([]byte("Allu"))
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

	now := time.Now()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	url.ID = fmt.Sprintf("%d", rnd.Uint32())[:6]
	url.CreatedAt = now
	url.UpdatedAt = now

	err = c.Repository.Create(&url, r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Controller(Create) error:", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"shortId\": \"%s\"}", url.ID)))
}

func (c *Controller) GetUrl(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	url, err := c.Repository.FindByID(id, r.Context())
	if err != nil {
		slog.Error("Controller(GetUrl) error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url.Url, http.StatusSeeOther)
}
