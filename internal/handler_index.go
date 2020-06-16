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

			h.serv.writeOKMessage("show nonAuthorizedIndex.html", "")

			h.serv.Respond(w, nil, "views/nonAuthorizedIndex.html")

			return
		}

		h.serv.writeErrorLog(err)

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	h.serv.writeOKMessage("show index.html", login)

	h.serv.Respond(w, login, "views/index.html")
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
