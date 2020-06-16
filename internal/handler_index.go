package internal

import (
	"errors"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

type indexHandler struct {
	serv *Server
}

func newIndexHandler(serv *Server) *indexHandler {

	return &indexHandler{serv: serv}

}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	message, code, user := h.serv.getInfoForRespond(r)

	h.serv.writeResponse(w, message, code, user)
}

func (s *Server) getUserFromClaimsFromCookie(r *http.Request) (*model.User, error) {

	c, err := r.Cookie("token")

	if err != nil {
		return nil, err
	}

	tokenString := c.Value

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid token")
	}

	user, err := s.DB.UserCol.FindByLogin(claims.Login)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Server) getInfoForRespond(r *http.Request) (string, int, *model.User) {
	user, err := s.getUserFromClaimsFromCookie(r)

	if err != nil {

		return err.Error() + "you're not authorized, try /createUser or /authorize", http.StatusBadRequest, user
	}

	return "you're authorized ", http.StatusOK, user
}
