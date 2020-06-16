package internal

import (
	db "github.com/anastasja-hunko/smptServer/internal/database"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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

	err := s.configureLogger()
	if err != nil {
		return err
	}

	s.configureRouter()

	err = s.configureDatabase()

	if err != nil {
		return err
	}

	s.writeOKMessage("starting server...", "")

	return http.ListenAndServe(s.config.Port, s.router)
}

func (s *Server) Disconnect() {

	err := s.DB.Close()

	if err != nil {
		s.Logger.Error("can't close db connection...")
	}

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

	s.router.Handle("/", indexHandler).Methods("GET")

	sendHandler := NewSendHandler(s)

	s.router.Handle("/sendMail", sendHandler)

	userHandler := NewUserHandler(s)

	s.router.Handle("/createUser", userHandler)

	s.router.Handle("/changePassword", userHandler)

	s.router.Handle("/showUsers", userHandler)

	s.router.Handle("/delete", userHandler)

	autorHandler := NewAutorHandler(s)

	s.router.Handle("/authorize", autorHandler)

	s.router.Handle("/logout", autorHandler)
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

	err := s.DB.LogCol.Create(logMessage)
	if err != nil {
		s.Logger.Error(err)
	}

}

func (s *Server) writeErrorLog(err error) {

	s.Logger.Error(err)

	logMessage := model.NewLog(err.Error(), "")

	s.writeLog(logMessage)
}

func (s *Server) writeOKMessage(message string, login string) {

	s.Logger.Info(message)

	logMessage := model.NewLog(message, login)

	s.writeLog(logMessage)
}
