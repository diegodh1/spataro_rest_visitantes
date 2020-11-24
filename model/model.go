package model

import (
	"time"
)

//Area structura para las diferentes areas de spataro
type Area struct {
	AreaID            string     `gorm:"primaryKey;"`
	AreaExtencion     string     `gorm:"default:'';"`
	AreaFechaCreacion *time.Time `gorm:"type:timestamptz;default:NOW();"`
	AreaEstado        bool       `gorm:"default:true;"`
}

//Documento structura para las diferentes documentos que puede tener una persona
type Documento struct {
	DocumentoID            string     `gorm:"primaryKey;"`
	DocumentoEstado        bool       `gorm:"default:true;"`
	DocumentoFechaCreacion *time.Time `gorm:"type:timestamptz;default:NOW();"`
}

//Empresa structura para las diferentes empresas que visitan a spataro
type Empresa struct {
	EmpresaID            string     `gorm:"primaryKey;"`
	EmpresaEstado        bool       `gorm:"default:true;"`
	EmpresaFechaCreacion *time.Time `gorm:"type:timestamptz;default:NOW();"`
}

//Permiso structura para las diferentes permisos que puede tener una persona
type Permiso struct {
	PermisoID            string `gorm:"primaryKey;"`
	PermisoDesc          string
	PermisoEstado        bool       `gorm:"default:true;"`
	PermisoFechaCreacion *time.Time `gorm:"type:timestamptz;default:NOW();"`
}

//Usuario structura para las diferentes usuarios que usaran la aplicacion
type Usuario struct {
	UsuarioID            uint64 `gorm:"primaryKey;"`
	UsuarioContrasena    string
	UsuarioNombre        string
	UsuarioApellido      string
	UsuarioEstado        bool       `gorm:"default:true;"`
	UsuarioFechaCreacion *time.Time `gorm:"type:timestamptz;default:NOW();"`
}

//UsuarioPermiso structura para los diferentes permisos que tiene un usuario
type UsuarioPermiso struct {
	PermisoID            string     `gorm:"primaryKey;"`
	UsuarioID            uint64     `gorm:"primaryKey;"`
	UsuarioPermisoEstado bool       `gorm:"default:true;"`
	UsuarioPermisoFecha  *time.Time `gorm:"type:timestamptz;default:NOW();"`
}

//Empleado structura para el empleado que trabaja en spataro
type Empleado struct {
	EmpleadoID            uint64 `gorm:"primaryKey;"`
	AreaID                string `gorm:"foreignKey:AreaID;"`
	EmpleadoNombre        string
	EmpleadoApellido      string
	EmpleadoEstado        bool       `gorm:"default:true;"`
	EmpleadoFechaCreacion *time.Time `gorm:"type:timestamptz;default:NOW();"`
}

//Visitante structura para las personas que visitan Spataro
type Visitante struct {
	VisitanteID            uint64 `gorm:"primaryKey;"`
	DocumentoID            string `gorm:"primaryKey;"`
	VisitanteNombre        string
	VisitanteApellido      string
	VisitanteCelular       string
	VisitanteCorreo        string
	VisitanteEstado        bool       `gorm:"default:true;"`
	VisitanteFechaCreacion *time.Time `gorm:"type:timestamptz;default:NOW();"`
}

//VisitanteDocumento structura para los diferentes documentos que tiene una persona
type VisitanteDocumento struct {
	VisitanteID             uint64 `gorm:"primaryKey;"`
	DocumentoID             string `gorm:"primaryKey;"`
	VisitanteDocNombre      string `gorm:"primaryKey;"`
	VisitanteDocReferencia  string
	VisitanteDocDescripcion string
	VisitanteDocFecha       *time.Time `gorm:"type:timestamptz;default:NOW();"`
	VisitanteDocPath        string
}

//VisitanteEmpresa structura para las visitas que hacen a spataro
type VisitanteEmpresa struct {
	VisitanteEmpresaRegistro uint64 `gorm:"primaryKey;"`
	VisitanteEmpresaHoras    uint64
	VisitanteID              uint64     `gorm:"foreignKey:VisitanteID;"`
	DocumentoID              string     `gorm:"foreignKey:DocumentoID;"`
	EmpresaID                string     `gorm:"foreignKey:EmpresaID;"`
	FechaEntrada             *time.Time `gorm:"type:timestamptz;default:NOW();"`
	FechaSalida              *time.Time `gorm:"type:timestamptz;"`
	FechaRealSalida          *time.Time `gorm:"type:timestamptz;"`
	Observaciones            string
	RegistroSalida           bool   `gorm:"default:false;"`
	EmpleadoID               uint64 `gorm:"foreignKey:EmpleadoID;"`
	UsuarioID                uint64 `gorm:"foreignKey:UsuarioID;"`
}
