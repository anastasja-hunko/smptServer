package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

//Сreate the JWT key used to create the signature
var jwtKey = []byte("nastya_key")

//Claims struct
type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

type autorHandler struct {
	serv *Server
}

func newAutorHandler(serv *Server) *autorHandler {

	return &autorHandler{serv: serv}
}

func (h *autorHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)

	defer cancel()

	if strings.Contains(r.URL.String(), "authorize") {

		h.authorizeHandler(ctx, rw, r)

	}

	if strings.Contains(r.URL.String(), "logout") {

		h.logout(ctx, rw)

	}
}

func (h *autorHandler) authorizeHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	var user = &model.User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {

		h.serv.writeResponse(ctx, rw, err.Error(), http.StatusBadRequest, nil)

		return

	}

	err = h.authorize(user, rw, ctx)
	if err != nil {

		h.serv.writeResponse(ctx, rw, err.Error(), http.StatusBadRequest, nil)

		return
	}

	h.serv.writeResponse(ctx, rw, "User was authorized", http.StatusOK, user)
}

/*If a user logs in with the correct credentials, this handler will
then set a cookie on the client side with the JWT value. Once a cookie is
set on a client, it is sent along with every request henceforth.
*/

func (h *autorHandler) authorize(ctx context.Context, u *model.User, rw http.ResponseWriter) error {

	user, err := h.serv.DB.UserCol.FindByLogin(ctx, u.Login)

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

	fmt.Println(tokenString)

	http.SetCookie(rw, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	return nil
}

//logout handler. If you're authorized, you see "quit" link on the index page.
//Results: Get: clean cookie and redirect
func (h *autorHandler) logout(ctx context.Context, rw http.ResponseWriter) {

	http.SetCookie(rw, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	h.serv.writeResponse(ctx, rw, "Logout", http.StatusOK, nil)
}
