package internal

import (
	"errors"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

//Сreate the JWT key used to create the signature
var jwtKey = []byte("nastya_key")

type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

type autorHandler struct {
	serv *Server
}

func NewAutorHandler(serv *Server) *autorHandler {
	return &autorHandler{serv: serv}
}

func (h *autorHandler) authorizeHandler(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		u := &model.User{
			Login:    r.FormValue("login"),
			Password: r.FormValue("password"),
		}

		err := h.authorize(u, rw)

		if err != nil {
			h.serv.writeErrorLog(err)

			rw.WriteHeader(http.StatusBadRequest)

			return
		}

		h.serv.writeOKMessage("User was authorized", u.Login)

		http.Redirect(rw, r, "/", 302)

		return
	}

	h.serv.Respond(rw, "Authorization", "views/userForm.html")
}

/*If a user logs in with the correct credentials, this handler will
then set a cookie on the client side with the JWT value. Once a cookie is
set on a client, it is sent along with every request henceforth.
*/

func (h *autorHandler) authorize(u *model.User, rw http.ResponseWriter) error {

	user, err := h.serv.DB.UserCol.FindByLogin(u.Login)

	if err != nil || !user.ComparePasswords(u.Password) {
		return errors.New("incorrect password or login")
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Login: user.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return err
	}

	http.SetCookie(rw, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	return nil
}

//logout handler. If you're authorized, you see "quit" link on the index page.
//Results: Get: clean cookie and redirect
func (h *autorHandler) Logout(rw http.ResponseWriter, r *http.Request) {

	http.SetCookie(rw, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	h.serv.writeOKMessage("Logout", "")

	http.Redirect(rw, r, "/", 302)
}
