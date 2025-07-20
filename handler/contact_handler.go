package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"porto/model"
	"porto/service"
	"strconv"

	"github.com/gorilla/mux"
)

type ContactHandler struct {
	Service service.ContactService
}

func NewContactHandler(s service.ContactService) *ContactHandler {
	return &ContactHandler{Service: s}
}

func (h *ContactHandler) GetContacts(w http.ResponseWriter, r *http.Request) {
	contacts, err := h.Service.GetAll(r.Context())
	if err != nil {
		log.Printf("[ContactHandler] GetContacts error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("[ContactHandler] GetContacts success, count: %d", len(contacts))
	json.NewEncoder(w).Encode(contacts)
}

func (h *ContactHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	var c model.Contact
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Printf("[ContactHandler] CreateContact decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Service.Create(r.Context(), &c); err != nil {
		log.Printf("[ContactHandler] CreateContact service error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	log.Printf("[ContactHandler] CreateContact success: %+v", c)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (h *ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[ContactHandler] DeleteContact invalid id: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Service.Delete(r.Context(), id); err != nil {
		log.Printf("[ContactHandler] DeleteContact service error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("[ContactHandler] DeleteContact success, id: %d", id)
	w.WriteHeader(http.StatusNoContent)
}

type ContactPageData struct {
	Address      string
	AddressDetail string
	Phone        string
	PhoneDesc    string
	Email        string
	EmailDesc    string
	Form         struct {
		Name    string
		Email   string
		Subject string
		Message string
	}
	FormMessage string
	Static      string
}

func (h *ContactHandler) RenderContactPage(w http.ResponseWriter, r *http.Request) {
	data := ContactPageData{
		Address:      "California, United States",
		AddressDetail: "Santa monica bullevard",
		Phone:        "00 (440) 9865 562",
		PhoneDesc:    "Mon to Fri 9am to 6 pm",
		Email:        "support@colorlib.com",
		EmailDesc:    "Send us your query anytime!",
	}
	// Ambil path static dari header, env, atau default
	staticPath := r.Header.Get("X-Static-Path")
	if staticPath == "" {
		staticPath = "/static"
	}
	data.Static = staticPath
	if r.Method == http.MethodPost {
		r.ParseForm()
		data.Form.Name = r.FormValue("name")
		data.Form.Email = r.FormValue("email")
		data.Form.Subject = r.FormValue("subject")
		data.Form.Message = r.FormValue("message")
		contact := model.Contact{
			Name:    data.Form.Name,
			Email:   data.Form.Email,
			Message: data.Form.Message,
		}
		err := h.Service.Create(r.Context(), &contact)
		if err != nil {
			data.FormMessage = err.Error()
		} else {
			data.FormMessage = "Pesan berhasil dikirim!"
			data.Form = struct {
				Name    string
				Email   string
				Subject string
				Message string
			}{} // reset form
		}
	}
	tmplPath := filepath.Join("WebView", "contact.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Template error: " + err.Error()))
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Render error: " + err.Error()))
	}
}
