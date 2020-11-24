package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	models "spataro/model"
	"strconv"

	"github.com/gorilla/mux"
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

//ALEXANDER METHODS

//ListUsers list all users
func ListUsers(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var users []models.Usuario
	db.Find(&users)
	if len(users) != 0 {
		respondJSON(w, http.StatusOK, users)
		return
	}
	respondError(w, http.StatusNotFound, "Not users found")
}

//UpdateUser update a user
func UpdateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := models.Usuario{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if queryRes := db.Omit("UsuarioContrasena", "UsuarioID", "UsuarioFechaCreacion").Updates(&user); queryRes.Error != nil || queryRes.RowsAffected == 0 {
		respondError(w, http.StatusBadRequest, "Error in operation or Not found")
		return
	}
	respondJSON(w, http.StatusOK, user)
}

//GetAllPermissions get all permissions
func GetAllPermissions(db *gorm.DB, w http.ResponseWriter) {
	var permisos []models.Permiso
	db.Find(&permisos)
	if len(permisos) != 0 {
		respondJSON(w, http.StatusOK, permisos)
		return
	}
	respondError(w, http.StatusNotFound, "Not found")
}

//UpdatePermissions update a permission
func UpdatePermissions(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	permiso := models.Permiso{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&permiso); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if queryRes := db.Updates(&permiso); queryRes.Error != nil || queryRes.RowsAffected == 0 {
		respondError(w, http.StatusBadRequest, "Error in operation or Not found")
		return
	}
	// return all the profiles
	var permisos []models.Permiso
	if err := db.Find(&permisos).Error; err != nil {
		respondError(w, http.StatusBadRequest, "Error in operation or Not found")
		return
	}
	respondJSON(w, http.StatusOK, permisos)
}

//GetUser get a user
func GetUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := models.Usuario{}
	user.UsuarioID, _ = strconv.ParseUint(mux.Vars(r)["UsuarioID"], 10, 64)

	if err := db.Find(&user).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, user)
}

//GetAllUserPermissions get all the permissions from a user
func GetAllUserPermissions(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var permisos []models.Permiso
	userPermission := models.UsuarioPermiso{}
	userPermission.UsuarioID, _ = strconv.ParseUint(mux.Vars(r)["UsuarioID"], 10, 64)

	rows, err := db.Raw("SELECT * FROM (SELECT permiso_id FROM usuario_permiso WHERE usuario_permiso_estado = true AND usuario_id = 10109182) AS foo NATURAL JOIN permiso",
		true, userPermission.UsuarioID).Rows()

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Data base error")
		return
	}

	defer rows.Close()
	entered := rows.Next()

	if !entered {
		respondJSON(w, http.StatusNoContent, permisos)
		return
	}

	for entered {
		if db.ScanRows(rows, &permisos) != nil {
			respondError(w, http.StatusInternalServerError, "Data base error")
			return
		}
		entered = rows.Next()
	}

	respondJSON(w, http.StatusOK, permisos)
}
