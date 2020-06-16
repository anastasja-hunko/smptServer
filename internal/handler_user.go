package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"net/http"
	"strings"
)

type userHandler struct {
	serv *Server
}

func NewUserHandler(s *Server) *userHandler {
	return &userHandler{serv: s}
}

func (h *userHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if strings.Contains(r.URL.String(), "createUser") {

		h.createUser(rw, r)

	}

	if strings.Contains(r.URL.String(), "changePassword") {

		h.changePassword(rw, r)

	}

	if strings.Contains(r.URL.String(), "showUsers") {

		h.showUsers(rw, r)

	}

	if strings.Contains(r.URL.String(), "delete") {

		h.deleteUser(rw, r)

	}
}

//create user handler. If you're not authorized, you see "create a user" link on the index page.
//Results: Get: show form for data input
//		   Post: create user and save it in db, and redirect to index page
func (h *userHandler) createUser(rw http.ResponseWriter, r *http.Request) {

	var user = &model.User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {

		h.serv.WriteResponse(rw, err.Error(), http.StatusBadRequest, nil)

		return

	}

	if user.Login == "" || user.Password == "" {

		h.serv.WriteResponse(rw, "Login or password are empty", http.StatusBadRequest, nil)

		return
	}

	err = h.registerUser(user)
	if err != nil {

		h.serv.WriteResponse(rw, err.Error(), http.StatusBadRequest, nil)

		return
	}

	h.serv.WriteResponse(rw, "user was created: "+user.Login, http.StatusCreated, user)

}

func (h *userHandler) registerUser(u *model.User) error {

	user, _ := h.serv.DB.UserCol.FindByLogin(u.Login)

	if user != nil {
		return errors.New("user's already existed with login:" + u.Login)
	}

	err := h.serv.DB.UserCol.Create(u)

	if err != nil {
		return err
	}

	return nil
}

//change password handler. If you're not authorized, you see "forgot a password" link on the index page.
//Results: Get: show form for data input
//		   Post: update user and save it in db, and redirect to index page
func (h *userHandler) changePassword(rw http.ResponseWriter, r *http.Request) {

	var user = &model.User{}

	_ = json.NewDecoder(r.Body).Decode(user)

	if user.Login == "" || user.Password == "" {

		h.serv.WriteResponse(rw, "Login or password are empty", http.StatusBadRequest, nil)

		return
	}

	err := h.serv.DB.UserCol.UpdatePassword(user)

	if err != nil {

		h.serv.WriteResponse(rw, err.Error(), http.StatusBadRequest, user)

		return
	}

	h.serv.WriteResponse(rw, "password was updated", http.StatusOK, user)

	return

}

//create user handler. If you're authorized, you see "delete an user" link on the index page.
//Results: Get: show all users in the table
func (h *userHandler) showUsers(rw http.ResponseWriter, r *http.Request) {

	users, err := h.serv.DB.UserCol.FindAll()

	if err != nil {

		h.serv.WriteResponse(rw, err.Error(), http.StatusBadRequest, nil)

		return
	}

	h.serv.WriteResponsePlus(rw, "users", http.StatusOK, nil, users)
}

//delete user handler. If you chose the link 'delete an user' in user list. See showUsers().
//Results: Get: delete the user from db
func (h *userHandler) deleteUser(rw http.ResponseWriter, r *http.Request) {

	login := fmt.Sprint(r.URL.Query().Get("login"))

	err := h.serv.DB.UserCol.UpdateActive(login)

	if err != nil {

		h.serv.WriteResponse(rw, err.Error(), http.StatusBadRequest, nil)

		return
	}

	user, _ := h.serv.getUserFromClaimsFromCookie(r)

	if login == user.Login {

		h.serv.WriteResponse(rw, "user active was updated and logged out", http.StatusOK, user)

	} else {

		h.serv.WriteResponse(rw, "user active was updated", http.StatusOK, nil)

	}
}
