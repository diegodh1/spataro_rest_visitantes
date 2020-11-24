package api

import (
	"fmt"
	"log"
	"net/http"
	"spataro/config"
	"spataro/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

//App struct
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

//Initialize initialize db
func (a *App) Initialize(config *config.Config) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		config.DB.Host,
		config.DB.Port,
		config.DB.Username,
		config.DB.Database,
		config.DB.Password,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		log.Fatal("Could not connect database")
	}
	a.DB = db
	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) setRouters() {
	a.Post("/login", a.UserLogin)
	a.Post("/createUser", a.CreateUser)
	a.Post("/createPermission", a.CreatePermission)
	a.Post("/assignPermission", a.AssignUserPermission)
	a.Post("/createGuest", a.CreateGuest)
	a.Post("/getGuest/", a.SearchGuest)
	a.Post("/createDocGuest/", a.CreateDocGuest)
	a.Post("/createGuestCompany/", a.CreateGuestCompany)
	a.Post("/finishGuestCompany/", a.UpdateVisit)
	a.Post("/updateGuest/", a.UpdateGuest)
	a.Post("/getAllDocumentsFromGuest/", a.GetGuestDocuments)
	a.Post("/getDocumentBase64/", a.GetDocumentBase64)
	a.Post("/getCompaniesGuest/", a.GetGuestCompanies)
	a.Post("/updateUser", a.UpdateUser)
	a.Post("/updatePermissions", a.UpdatePermissions)
	//GET
	a.Get("/getAllEmployees/", a.GetAllEmployees)
	a.Get("/getAllDocuments/", a.GetAllDocuments)
	a.Get("/listUsers", a.ListUsers)
	a.Get("/getAllPermissions", a.GetAllPermissions)
	a.Get("/getUser/{UsuarioID}", a.GetUser)
}

//Get all get functions
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

//Post all Post functions
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

//Put all Put functions
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

//UserLogin let the user log in the app
func (a *App) UserLogin(w http.ResponseWriter, r *http.Request) {
	handler.UserLogin(a.DB, w, r)
}

//CreateUser lets create a new user
func (a *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	handler.CreateUser(a.DB, w, r)
}

//CreatePermission lets create a new 'permiso'
func (a *App) CreatePermission(w http.ResponseWriter, r *http.Request) {
	handler.CreatePermission(a.DB, w, r)
}

//AssignUserPermission lets assign a 'permiso' to a user
func (a *App) AssignUserPermission(w http.ResponseWriter, r *http.Request) {
	handler.AssignUserPermission(a.DB, w, r)
}

//CreateGuest creates a new guest in the database
func (a *App) CreateGuest(w http.ResponseWriter, r *http.Request) {
	handler.CreateGuest(a.DB, w, r)
}

//SearchGuest search a guest in the database
func (a *App) SearchGuest(w http.ResponseWriter, r *http.Request) {
	handler.SearchGuest(a.DB, w, r)
}

//CreateDocGuest creates a new guest in the database
func (a *App) CreateDocGuest(w http.ResponseWriter, r *http.Request) {
	handler.CreateDocGuest(a.DB, w, r)
}

//CreateGuestCompany creates a new guest in the database
func (a *App) CreateGuestCompany(w http.ResponseWriter, r *http.Request) {
	handler.CreateGuestCompany(a.DB, w, r)
}

//UpdateVisit finish the visit
func (a *App) UpdateVisit(w http.ResponseWriter, r *http.Request) {
	handler.UpdateVisit(a.DB, w, r)
}

//GetAllDocuments this function returns all the documents type
func (a *App) GetAllDocuments(w http.ResponseWriter, r *http.Request) {
	handler.GetAllDocuments(a.DB, w, r)
}

//UpdateGuest this function updates a guest
func (a *App) UpdateGuest(w http.ResponseWriter, r *http.Request) {
	handler.UpdateGuest(a.DB, w, r)
}

//GetGuestDocuments Get all the documents from a guest
func (a *App) GetGuestDocuments(w http.ResponseWriter, r *http.Request) {
	handler.GetGuestDocuments(a.DB, w, r)
}

//GetDocumentBase64 Get a document in base 64
func (a *App) GetDocumentBase64(w http.ResponseWriter, r *http.Request) {
	handler.GetDocumentBase64(a.DB, w, r)
}

//GetGuestCompanies Get all the companies that were visited by the guest
func (a *App) GetGuestCompanies(w http.ResponseWriter, r *http.Request) {
	handler.GetGuestCompanies(a.DB, w, r)
}

//GetAllEmployees Get all the employees from the database
func (a *App) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	handler.GetAllEmployees(a.DB, w, r)
}

// FUNCIONES ALEXANDER
func (a *App) ListUsers(w http.ResponseWriter, r *http.Request) {
	handler.ListUsers(a.DB, w, r)
}

func (a *App) UpdateUser(w http.ResponseWriter, r *http.Request) {
	handler.UpdateUser(a.DB, w, r)
}

func (a *App) GetAllPermissions(w http.ResponseWriter, _ *http.Request) {
	handler.GetAllPermissions(a.DB, w)
}

func (a *App) UpdatePermissions(w http.ResponseWriter, r *http.Request) {
	handler.UpdatePermissions(a.DB, w, r)
}

func (a *App) GetUser(w http.ResponseWriter, r *http.Request) {
	handler.GetUser(a.DB, w, r)
}

//Run run app
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(a.Router)))
}
