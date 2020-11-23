package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	models "spataro/model"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

//GetAllDocuments this function returns all the documents type
func GetAllDocuments(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	documents := []models.Documento{}
	db.Where("documento_estado = ?", true).Find(&documents)
	respondJSON(w, http.StatusOK, JSONResponse{Payload: documents, Message: "Documentos encontrados!!"})
}

//CreateGuest this function allows to create a new 'visitante'
func CreateGuest(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	guest := models.Visitante{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&guest); err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Petición mal estructurada"})
		fmt.Println(err.Error())
		return
	}
	defer r.Body.Close()
	guestTemp := getGuestOrNull(db, guest.VisitanteID, guest.DocumentoID)
	if guestTemp != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "el visitante ya está registrado"})
		return
	}

	if err := db.Omit("VisitanteFechaCreacion").Save(&guest).Error; err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Por favor verificar los campos obligatorios"})
		return
	}
	respondJSON(w, http.StatusOK, JSONResponse{Message: "Visitante creado con éxito!!"})
}

//UpdateVisit this function allows to create a new 'visitante'
func UpdateVisit(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	guest := models.VisitanteEmpresa{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&guest); err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Petición mal estructurada"})
		fmt.Println(err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Model(&guest).Select("RegistroSalida", "Observaciones", "FechaRealSalida").Updates(guest).Error; err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Por favor verificar los campos obligatorios"})
		return
	}
	respondJSON(w, http.StatusOK, JSONResponse{Message: "Visita Finalizada!!"})
}

//CreateCompany create a new company in the database
func CreateCompany(db gorm.DB, empresaID string) (string, error) {
	company := models.Empresa{EmpresaID: empresaID}

	if err := db.Omit("EmpresaFechaCreacion").Save(&company).Error; err != nil {
		return "", errors.New(err.Error())
	}
	return "empresa creada con éxito", nil
}

//CreateGuestCompany create a new visit in the database
func CreateGuestCompany(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	guestCompany := models.VisitanteEmpresa{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&guestCompany); err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Petición mal estructurada"})
		fmt.Println(err.Error())
		return
	}
	defer r.Body.Close()
	guestCompany.EmpresaID = strings.ToUpper(guestCompany.EmpresaID)
	if guestCompany.EmpresaID == "" {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "La empresa no puede estar vacia"})
		return
	}
	_, err := CreateCompany(*db, guestCompany.EmpresaID)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: err.Error()})
		return
	}
	if err := db.Create(&guestCompany).Error; err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Por favor verificar los campos obligatorios"})
		return
	}
	respondJSON(w, http.StatusOK, JSONResponse{Message: "Visita creada con éxito!!"})
}

//SearchGuest this function search a 'visitante' in the database
func SearchGuest(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	guest := models.Visitante{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&guest); err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Petición mal estructurada"})
		fmt.Println(err.Error())
		return
	}
	defer r.Body.Close()
	guestTemp := getGuestOrNull(db, guest.VisitanteID, guest.DocumentoID)
	if guestTemp != nil {
		respondJSON(w, http.StatusOK, JSONResponse{Payload: guestTemp, Message: "registro encontrado"})
		return
	}
	respondJSON(w, http.StatusOK, JSONResponse{Message: "Visitante no está registrado"})
}

//CreateDocGuest this function creates a new doc that belongs to a guest
func CreateDocGuest(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	file, handler, _ := r.FormFile("file")
	visitanteAux := r.FormValue("visitanteID")
	visitanteID, err := strconv.Atoi(visitanteAux)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "El id de usuario debe ser numérico"})
		return
	}
	docGuest := models.VisitanteDocumento{}
	docGuest.VisitanteID = uint64(visitanteID)
	docGuest.DocumentoID = r.FormValue("documentoID")
	docGuest.VisitanteDocNombre = strings.ToUpper(r.FormValue("docNombre"))
	docGuest.VisitanteDocReferencia = strings.ToUpper(r.FormValue("docReferencia"))
	docGuest.VisitanteDocDescripcion = r.FormValue("docDescripcion")

	if err := db.Omit("VisitanteDocFecha").Save(&docGuest).Error; err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Por favor verificar los campos obligatorios"})
		return
	}
	if file != nil {
		path := `C:\visitas\` + visitanteAux
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.MkdirAll(path, 0755)
		}
		destination, err := os.Create(path + `\` + handler.Filename)
		if err != nil {
			respondJSON(w, http.StatusOK, JSONResponse{Message: "registro realizado sin archivo adjunto!!"})
			return
		}
		defer destination.Close()
		_, _ = io.Copy(destination, file)
		docGuest.VisitanteDocPath = path + `\` + handler.Filename
		file.Close()
		if err := db.Omit("VisitanteDocFecha").Save(&docGuest).Error; err != nil {
			respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Por favor verificar los campos obligatorios"})
			return
		}
		respondJSON(w, http.StatusOK, JSONResponse{Message: "registro realizado con archivo adjunto!!"})
		return
	}
	respondJSON(w, http.StatusOK, JSONResponse{Message: "registro realizado sin archivo adjunto!!"})
}

// get a user whose ID is equal to the param
func getGuestOrNull(db *gorm.DB, visitanteID uint64, documentoID string) *models.Visitante {
	guest := models.Visitante{}
	if err := db.Where("visitante_estado = ?", true).First(&guest, models.Visitante{VisitanteID: visitanteID, DocumentoID: documentoID}).Error; err != nil {
		return nil
	}
	return &guest
}
