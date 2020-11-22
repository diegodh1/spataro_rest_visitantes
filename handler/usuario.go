package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	models "spataro/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//CreateUser this function allows to create a new user
func CreateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := models.Usuario{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&user); err != nil {
		respondJSON(w, http.StatusUnauthorized, JSONResponse{Message: "Error interno del servidor"})
		return
	}
	defer r.Body.Close()
	if user.UsuarioContrasena == "" {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Debe ingresar una contraseña"})
		return
	}
	userTemp := getUserOrNull(db, user.UsuarioID)
	if userTemp != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "el usuario ya está registrado"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.UsuarioContrasena), 10)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, JSONResponse{Message: "Error Interno del servidor"})
		return
	}
	user.UsuarioContrasena = bytes.NewBuffer(hash).String()

	if err := db.Create(&user).Error; err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: err.Error()})
	}
	respondJSON(w, http.StatusOK, JSONResponse{Message: "Usuario creado con éxito!!"})
}
