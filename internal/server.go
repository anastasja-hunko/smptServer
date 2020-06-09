package internal

import (
	db "github.com/anastasja-hunko/smptServer/internal/database"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"html/template"
	"net/http"
)

type Server struct {
	config *Config
	Logger *logrus.Logger
	router *mux.Router
	DB     *db.Database
}

func New(config *Config) *Server {

	return &Server{
		config: config,
		Logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

//start a server
func (s *Server) Start() error {

	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	if err := s.configureDatabase(); err != nil {
		return err
	}

	defer func() {
		if err := s.DB.Close(); err != nil {
			s.Logger.Error("can't close db connection...")
		}
	}()

	s.writeOKMessage("starting server...", "")

	return http.ListenAndServe(s.config.Port, s.router)
}

func (s *Server) configureLogger() error {

	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.Logger.SetLevel(level)
	return nil
}

//endpoints
func (s *Server) configureRouter() {

	indexHandler := NewIndexHandler(s)

	s.router.HandleFunc("/", indexHandler.indexPage).Methods("GET")

	sendHandler := NewSendHandler(s)

	s.router.HandleFunc("/sendMail", sendHandler.sendHandler)

	userHandler := NewUserHandler(s)

	s.router.HandleFunc("/createUser", userHandler.CreateUser)

	s.router.HandleFunc("/changePassword", userHandler.changePassword)

	s.router.HandleFunc("/showUsers", userHandler.showUsers)

	s.router.HandleFunc("/delete", userHandler.deleteUser)

	autorHandler := NewAutorHandler(s)

	s.router.HandleFunc("/authorize", autorHandler.authorizeHandler)

	s.router.HandleFunc("/logout", autorHandler.Logout)
}

//configure database with config
func (s *Server) configureDatabase() error {

	dbase := db.New(s.config.DbConfig)

	err := dbase.Open()

	if err != nil {
		return err
	}

	s.DB = dbase
	return nil
}

//execute html template
func executeTemplate(page string, w http.ResponseWriter, data interface{}) error {

	tmpl := template.Must(template.ParseFiles(page))

	return tmpl.Execute(w, data)
}

//action's respond when everything is OK
func (s *Server) Respond(rw http.ResponseWriter, data interface{}, page string) {

	err := executeTemplate(page, rw, data)
	if err != nil {
		s.writeErrorLog(err)
	}
}

func (s *Server) writeLog(logMessage *model.Log) {
	done := make(chan bool)

	go func() {
		err := s.DB.Log().Create(logMessage, done)
		if err != nil {
			s.Logger.Error(err)
		}
	}()

	fmt.Println(<-done)
}

func (s *Server) writeErrorLog(err error) {

	s.Logger.Error(err)

	logMessage := model.NewLog(err.Error(), "")

	s.writeLog(logMessage)
}

func (s *Server) writeTextErrorLog(err string) {

	s.Logger.Error(err)

	logMessage := model.NewLog(err, "")

	s.writeLog(logMessage)
}

func (s *Server) writeOKMessage(message string, login string) {

	s.Logger.Info(message)

	logMessage := model.NewLog(message, login)

	s.writeLog(logMessage)
}
