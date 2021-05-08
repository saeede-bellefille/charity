package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"SABKAD/pkg/model"
)

type Handler struct {
	creator func() model.Model
	db      *gorm.DB
	name    string
}

func New(creator func() model.Model, db *gorm.DB, name string) *Handler {
	return &Handler{creator, db, name}
}

// Create creates an object in database
func (h *Handler) Create(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	obj := h.creator()

	if err := json.NewDecoder(r.Body).Decode(obj); err != nil && err != io.EOF {
		log.Println("[ERROR] deserializing", h.name, err)
		http.Error(rw, fmt.Sprintf("Error reading %s", h.name), http.StatusBadRequest)
		return
	}

	objects, err := obj.Find(h.db)
	if err != nil {
		log.Println("Unable to scan the row", h.name, err)
		http.Error(rw, fmt.Sprintf("Error finding %s", h.name), http.StatusBadRequest)
		return
	}
	if len(objects) != 0 {
		fmt.Println("[ERROR] already exists this", h.name)
		http.Error(rw, fmt.Sprintf("already exists %s", h.name), http.StatusNotAcceptable)
		return
	}

	obj.Initialize(h.db)
	if err := obj.Validate(); err != nil {
		log.Println("[ERROR] validation", h.name, err)
		http.Error(rw, fmt.Sprintf("ERROR validation: %s", err), http.StatusBadRequest)
		return
	}

	if err := h.db.Create(obj).Error; err != nil {
		log.Println("Cant insert new", h.name, err)
		http.Error(rw, fmt.Sprintf("Cant insert new %s: %v", h.name, err), http.StatusNotAcceptable)
		return
	}

	if err := json.NewEncoder(rw).Encode(obj); err != nil {
		log.Println("Unable to marshal json", h.name, err)
		http.Error(rw, fmt.Sprintf("Unable to marshal json: %v", err), http.StatusInternalServerError)
	}
}

//Search finds array of objects in database
func (h *Handler) Search(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	query := r.URL.Query()
	obj := h.creator()
	obj.Load(query)

	objects, err := obj.Find(h.db)
	if err != nil {
		log.Println("[ERROR] Unable to scan the row", h.name, err)
		http.Error(rw, fmt.Sprintf("ERROR Unable to scan the row: %s", err), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(rw).Encode(objects); err != nil {
		log.Println("[ERROR] Unable to marshal json", h.name, err)
		http.Error(rw, fmt.Sprintf("Unable to marshal json: %v", err), http.StatusInternalServerError)
	}
}

// Update update object details in the database
func (h *Handler) Update(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "PUT")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	obj := h.creator()
	obj.Load(model.MapGetter(mux.Vars(r)))

	objects, err := obj.Find(h.db)
	if err != nil {
		log.Println("Unable to scan the row", h.name, err)
		http.Error(rw, fmt.Sprintf("ERROR Unable to scan the row: %s", err), http.StatusInternalServerError)
		return
	}
	if len(objects) == 0 {
		fmt.Println("[ERROR] not exists this database", h.name)
		http.Error(rw, "ERROR not exists this database", http.StatusNotFound)
		return
	}

	obj = objects[0]
	if err := json.NewDecoder(r.Body).Decode(obj); err != nil && err != io.EOF {
		log.Println("[ERROR] deserializing", h.name, err)
		http.Error(rw, fmt.Sprintf("Error reading %s: %s", h.name, err), http.StatusBadRequest)
		return
	}
	obj.Initialize(h.db)

	result := h.db.Save(obj)
	if result.Error != nil {
		log.Println("[ERROR] cannot save", h.name, result.Error)
		http.Error(rw, fmt.Sprintf("Error saving %s: %s", h.name, result.Error), http.StatusInternalServerError)
		return
	}

	fmt.Printf("object updated successfully. Total rows/record affected %d", result.RowsAffected)
	if err := json.NewEncoder(rw).Encode(obj); err != nil {
		log.Println("[ERROR] Unable to marshal json", h.name, err)
		http.Error(rw, fmt.Sprintf("Unable to marshal json: %v", err), http.StatusInternalServerError)
	}
}

// Delete removes detail in the database
func (h *Handler) Delete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "DELETE")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	obj := h.creator()
	obj.Load(model.MapGetter(mux.Vars(r)))

	result := h.db.Delete(obj)
	if result.Error != nil {
		log.Println("[ERROR] cannot delete", h.name, result.Error)
		http.Error(rw, fmt.Sprintf("Error deleting %s: %s", h.name, result.Error), http.StatusInternalServerError)
		return
	}
	fmt.Printf("object delete successfully. Total rows/record affected %d", result.RowsAffected)
	rw.WriteHeader(http.StatusNoContent)
}
