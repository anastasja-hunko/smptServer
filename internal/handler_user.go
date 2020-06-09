package internal

import (
	"errors"
	"fmt"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"net/http"
)

type userHandler struct {
	serv *Server
}

func NewUserHandler(s *Server) *userHandler {
	return &userHandler{serv: s}
}

//create user handler. If you're not authorized, you see "create a user" link on the index page.
//Results: Get: show form for data input
//		   Post: create user and save it in db, and redirect to index page
func (h *userHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		user := &model.User{
			Login:    r.FormValue("login"),
			Password: r.FormValue("password"),
		}

		err := h.registerUser(user)

		if err != nil {
			h.serv.writeErrorLog(err)

			rw.WriteHeader(http.StatusBadRequest)

			return
		}
		h.serv.writeOKMessage("user was created: "+user.Login, "")

		http.Redirect(rw, r, "/", 302)

		return
	}

	h.serv.writeOKMessage("show userForm.html for creating user", "")

	h.serv.Respond(rw, "Create a user", "views/userForm.html")

}

func (h *userHandler) registerUser(u *model.User) error {

	user, _ := h.serv.DB.User().FindByLogin(u.Login)

	if user != nil {
		return errors.New("user's already existed with login:" + u.Login)
	}

	err := h.serv.DB.User().Create(u)

	if err != nil {
		return err
	}

	return nil
}

//change password handler. If you're not authorized, you see "forgot a password" link on the index page.
//Results: Get: show form for data input
//		   Post: update user and save it in db, and redirect to index page
func (h *userHandler) changePassword(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		user := &model.User{
			Login:    r.FormValue("login"),
			Password: r.FormValue("password"),
		}

		err := h.serv.DB.User().Update(user)

		if err != nil {

			h.serv.writeErrorLog(err)

			w.WriteHeader(http.StatusBadRequest)

			return
		}
		h.serv.writeOKMessage("password was updated: "+user.Login, "")

		http.Redirect(w, r, "/", 302)

		return
	}

	h.serv.writeOKMessage("show userForm.html for creating user", "")

	h.serv.Respond(w, "Change password", "views/userForm.html")
}

//create user handler. If you're authorized, you see "delete an user" link on the index page.
//Results: Get: show all users in the table
func (h *userHandler) showUsers(w http.ResponseWriter, r *http.Request) {

	users, err := h.serv.DB.User().FindAll()

	if err != nil {

		h.serv.writeErrorLog(err)

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	h.serv.Respond(w, users, "views/users.html")

	return
}

//delete user handler. If you chose the link 'delete an user' in user list. See showUsers().
//Results: Get: delete the user from db
func (h *userHandler) deleteUser(w http.ResponseWriter, r *http.Request) {

	login := fmt.Sprint(r.URL.Query().Get("login"))

	err := h.serv.DB.User().Delete(login)

	if err != nil {

		h.serv.writeErrorLog(err)

		return
	}

	cookieLogin, _ := getLoginFromClaimsFromCookie(r)

	if login == cookieLogin {

		http.Redirect(w, r, "/logout", 302)

	} else {

		http.Redirect(w, r, "/", 302)

	}
}
