package handlers

import (
	"gorm.io/gorm"
	"net/http"
)

type handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) handler {
	return handler{db}
}
func (h handler) GetAllFios(w http.ResponseWriter, r *http.Request) {}
func (h handler) GetFio(w http.ResponseWriter, r *http.Request)     {}
func (h handler) AddFio(w http.ResponseWriter, r *http.Request)     {}
func (h handler) UpdateFio(w http.ResponseWriter, r *http.Request)  {}
func (h handler) DeleteFio(w http.ResponseWriter, r *http.Request)  {}
