package api

import (
	"fmt"
	"log"
	"net/http"
	"spataro/config"
	"spataro/handler"

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

//Run run app
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
