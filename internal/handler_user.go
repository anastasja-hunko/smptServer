package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"net/http"
	"strings"
	"time"
)

type userHandler struct {
	serv *Server
}

func newUserHandler(s *Server) *userHandler {
	return &userHandler{serv: s}
}

func (h *userHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)

	defer cancel()

	if strings.Contains(r.URL.String(), "createUser") {

		h.createUser(ctx, rw, r)

	}

	if strings.Contains(r.URL.String(), "changePassword") {

		h.changePassword(ctx, rw, r)

	}

	if strings.Contains(r.URL.String(), "showUsers") {

		h.showUsers(ctx, rw)

	}

	if strings.Contains(r.URL.String(), "delete") {

		h.deleteUser(ctx, rw, r)

	}
}

//create user handler. If you're not authorized, you see "create a user" link on the index page.
//Results: Get: show form for data input
//		   Post: create user and save it in db, and redirect to index page
func (h *userHandler) createUser(ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	var user = &model.User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {

		h.serv.writeResponse(ctx, rw, err.Error(), http.StatusBadRequest, nil)

		return

	}

	if user.Login == "" || user.Password == "" {

		h.serv.writeResponse(ctx, rw, "Login or password are empty", http.StatusBadRequest, nil)

		return
	}

	err = h.registerUser(ctx, user)
	if err != nil {

		h.serv.writeResponse(ctx, rw, err.Error(), http.StatusBadRequest, nil)

		return
	}

	h.serv.writeResponse(ctx, rw, "user was created: "+user.Login, http.StatusCreated, user)

}

func (h *userHandler) registerUser(ctx context.Context, u *model.User) error {

	user, _ := h.serv.DB.UserCol.FindByLogin(ctx, u.Login)

	if user != nil {
		return errors.New("user's already existed with login:" + u.Login)
	}

	err := h.serv.DB.UserCol.Create(ctx, u)
	if err != nil {
		return err
	}

	return nil
}

//change password handler. If you're not authorized, you see "forgot a password" link on the index page.
//Results: Get: show form for data input
//		   Post: update user and save it in db, and redirect to index page
func (h *userHandler) changePassword(ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	var user = &model.User{}

	_ = json.NewDecoder(r.Body).Decode(user)

	if user.Login == "" || user.Password == "" {

		h.serv.writeResponse(ctx, rw, "Login or password are empty", http.StatusBadRequest, nil)

		return
	}

	err := h.serv.DB.UserCol.UpdatePassword(ctx, user)

	if err != nil {

		h.serv.writeResponse(ctx, rw, err.Error(), http.StatusBadRequest, user)

		return
	}

	h.serv.writeResponse(ctx, rw, "password was updated", http.StatusOK, user)

	return

}

//create user handler. If you're authorized, you see "delete an user" link on the index page.
//Results: Get: show all users in the table
func (h *userHandler) showUsers(ctx context.Context, rw http.ResponseWriter) {

	users, err := h.serv.DB.UserCol.FindAll(ctx)

	if err != nil {

		h.serv.writeResponse(ctx, rw, err.Error(), http.StatusBadRequest, nil)

		return
	}

	h.serv.writeResponsePlus(ctx, rw, "users", http.StatusOK, nil, users)
}

//delete user handler. If you chose the link 'delete an user' in user list. See showUsers().
//Results: Get: delete the user from db
func (h *userHandler) deleteUser(ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	login := fmt.Sprint(r.URL.Query().Get("login"))

	err := h.serv.DB.UserCol.UpdateActive(ctx, login)

	if err != nil {

		h.serv.writeResponse(ctx, rw, err.Error(), http.StatusBadRequest, nil)

		return
	}

	user, _ := h.serv.getUserFromClaimsFromCookie(ctx, r)

	if login == user.Login {

		h.serv.writeResponse(ctx, rw, "user active was updated and logged out", http.StatusOK, user)

	} else {

		h.serv.writeResponse(ctx, rw, "user active was updated", http.StatusOK, nil)

	}
}
