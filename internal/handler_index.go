package internal

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

type indexHandler struct {
	serv *Server
}

func NewIndexHandler(serv *Server) *indexHandler {

	return &indexHandler{serv: serv}

}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	login, err := getLoginFromClaimsFromCookie(r)

	if err != nil {

		if err == http.ErrNoCookie || err == jwt.ErrSignatureInvalid || err.Error() == "Invalid token" {

			h.serv.WriteResponse(w, "you're not authorized, try /createUser or /authorize", http.StatusOK, nil)

			return
		}

		h.serv.WriteResponse(w, err.Error(), http.StatusBadRequest, nil)

		return
	}

	h.serv.WriteResponse(w, "user authorized "+login, http.StatusOK, nil)
}

func getLoginFromClaimsFromCookie(r *http.Request) (string, error) {

	c, err := r.Cookie("token")

	if err != nil {
		return "", err
	}

	tokenString := c.Value

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("Invalid token")
	}

	return claims.Login, nil
}
