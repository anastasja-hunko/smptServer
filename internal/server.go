package internal

import (
	"context"
	"encoding/json"
	db "github.com/anastasja-hunko/smptServer/internal/database"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

//Server struct
type Server struct {
	Config  *Config
	Logger  *logrus.Logger
	router  *mux.Router
	DB      *db.Database
	context *context.Context
}

//New - returns initialized server
func New(config *Config) *Server {

	return &Server{
		Config: config,
		Logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

//Start a server
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

	s.Logger.Print("starting server...")

	return http.ListenAndServe(s.Config.Port, s.router)
}

//Restart a server
func (s *Server) Restart() error {

	err := s.configureDatabase()
	if err != nil {
		return err
	}

	s.Logger.Print("restarting server...")

	return http.ListenAndServe(s.Config.Port, s.router)
}

//Disconnect the server
func (s *Server) Disconnect() {

	err := s.DB.Close()

	if err != nil {
		s.Logger.Error("can't close db connection...")
	}

}

func (s *Server) configureLogger() error {

	level, err := logrus.ParseLevel(s.Config.LogLevel)
	if err != nil {
		return err
	}

	s.Logger.SetLevel(level)
	return nil
}

//endpoints
func (s *Server) configureRouter() {

	indexHandler := newIndexHandler(s)

	s.router.Handle("/", indexHandler).Methods("GET")

	sendHandler := newSendHandler(s)

	s.router.Handle("/sendMail", sendHandler)

	userHandler := newUserHandler(s)

	s.router.Handle("/createUser", userHandler)

	s.router.Handle("/changePassword", userHandler)

	s.router.Handle("/showUsers", userHandler)

	s.router.Handle("/delete", userHandler)

	autorHandler := newAutorHandler(s)

	s.router.Handle("/authorize", autorHandler)

	s.router.Handle("/logout", autorHandler)
}

//configure database with config
func (s *Server) configureDatabase() error {

	dbase := db.New(s.Config.DbConfig)

	err := dbase.Open()

	if err != nil {
		return err
	}

	s.DB = dbase
	return nil
}

func (s *Server) writeLog(ctx context.Context, logMessage string, user *model.User) {

	s.Logger.Error(logMessage)

	if user != nil {

		err := s.DB.UserCol.UpdateUserLog(ctx, user, logMessage)
		if err != nil {
			s.Logger.Error(err)
		}

		return

	}

	message := model.NewLog(logMessage)

	err := s.DB.LogCol.Create(ctx, message)
	if err != nil {
		s.Logger.Error(err)
	}
}

func (s *Server) writeResponse(
	ctx context.Context,
	w http.ResponseWriter,
	message string,
	status int,
	user *model.User,
) {

	s.writeResponsePlus(ctx, w, message, status, user, nil)
}

func (s *Server) writeResponsePlus(
	ctx context.Context,
	w http.ResponseWriter,
	message string,
	status int,
	user *model.User,
	addedData interface{},
) {

	s.writeLog(ctx, message, user)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	data := struct {
		Message   string      `json:"message"`
		AddedData interface{} `json:"data"`
	}{
		Message:   message,
		AddedData: addedData,
	}

	_ = json.NewEncoder(w).Encode(&data)
}
