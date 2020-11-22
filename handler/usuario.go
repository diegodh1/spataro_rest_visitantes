package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Petición mal estructurada"})
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

	if err := db.Omit("UsuarioFechaCreacion").Save(&user).Error; err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: err.Error()})
	}
	respondJSON(w, http.StatusOK, JSONResponse{Message: "Usuario creado con éxito!!"})
}

//CreatePermission this function allows to create a new 'permiso'
func CreatePermission(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	permission := models.Permiso{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&permission); err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Petición mal estructurada"})
		return
	}
	defer r.Body.Close()
	permissionTemp := getPermissionOrNull(db, permission.PermisoID)
	if permissionTemp != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "el permiso ya está registrado"})
		return
	}

	if err := db.Omit("PermisoFechaCreacion").Save(&permission).Error; err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: err.Error()})
	}
	respondJSON(w, http.StatusOK, JSONResponse{Message: "Perfil creado con éxito!!"})
}

//AssignUserPermission this function allows to create a new 'permiso'
func AssignUserPermission(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	userPermission := models.UsuarioPermiso{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&userPermission); err != nil {
		fmt.Println(err.Error())
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Petición mal estructurada"})
		return
	}
	defer r.Body.Close()
	if err := db.Omit("UsuarioPermisoFecha").Save(&userPermission).Error; err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: err.Error()})
		return
	}
	respondJSON(w, http.StatusOK, JSONResponse{Message: "Operación realizada con éxito!!"})
}

// get a user whose ID is equal to the param
func getPermissionOrNull(db *gorm.DB, permisoID string) *models.Permiso {
	user := models.Permiso{}
	if err := db.Where("permiso_estado = ?", true).First(&user, models.Permiso{PermisoID: permisoID}).Error; err != nil {
		return nil
	}
	return &user
}

// get a user whose ID is equal to the param
func getUserPermissionOrNull(db *gorm.DB, permisoID string, usuarioID uint64) *models.UsuarioPermiso {
	userPermiso := models.UsuarioPermiso{}
	if err := db.Where("permiso_id = ? AND usuario_id", permisoID, usuarioID).First(&userPermiso, models.UsuarioPermiso{PermisoID: permisoID, UsuarioID: usuarioID}).Error; err != nil {
		return nil
	}
	return &userPermiso
}
