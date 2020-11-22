package handler

import (
	"encoding/json"
	"net/http"
	models "spataro/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//UserLogin this function allow the user login if the user an password are correct
func UserLogin(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := models.Usuario{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&user); err != nil {
		respondJSON(w, http.StatusUnauthorized, JSONResponse{Message: "Error interno del servidor"})
		return
	}
	userTemp := getUserOrNull(db, user.UsuarioID)
	if userTemp == nil {
		respondJSON(w, http.StatusUnauthorized, JSONResponse{Message: "el usuario no está registrado ó está inactivo"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userTemp.UsuarioContrasena), []byte(user.UsuarioContrasena)); err != nil {
		respondJSON(w, http.StatusUnauthorized, JSONResponse{Message: "usuario y/o contraseña incorrecta"})
		return
	}
	defer r.Body.Close()
	profiles := getProfilesUser(db, userTemp.UsuarioID)
	anonymousStruct := struct {
		User     models.Usuario
		Perfiles *[]models.UsuarioPermiso
	}{*userTemp, &profiles}
	respondJSON(w, http.StatusOK, JSONResponse{Payload: anonymousStruct, Message: "Ingreso Realizado!"})
}

// this function get all the profiles that belongs to a user
func getProfilesUser(db *gorm.DB, userID uint64) []models.UsuarioPermiso {
	profiles := []models.UsuarioPermiso{}
	if err := db.Debug().Where("usuario_permiso_estado = ? AND usuario_id = ?", true, userID).Find(&profiles).Error; err != nil {
		return profiles
	}
	return profiles
}

// get a user whose ID is equal to the param
func getUserOrNull(db *gorm.DB, usuarioID uint64) *models.Usuario {
	user := models.Usuario{}
	if err := db.Where("usuario_estado = ?", true).First(&user, models.Usuario{UsuarioID: usuarioID}).Error; err != nil {
		return nil
	}
	return &user
}
