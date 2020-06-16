package internal

import (
	"context"
	"errors"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

type indexHandler struct {
	serv *Server
}

func newIndexHandler(serv *Server) *indexHandler {

	return &indexHandler{serv: serv}

}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)

	defer cancel()

	message, code, user := h.serv.getInfoForRespond(ctx, r)

	h.serv.writeResponse(ctx, w, message, code, user)
}

func (s *Server) getUserFromClaimsFromCookie(ctx context.Context, r *http.Request) (*model.User, error) {

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
		return nil, errors.New("invalid token")
	}

	user, err := s.DB.UserCol.FindByLogin(ctx, claims.Login)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Server) getInfoForRespond(ctx context.Context, r *http.Request) (string, int, *model.User) {

	user, err := s.getUserFromClaimsFromCookie(ctx, r)
	if err != nil {

		return err.Error() + "you're not authorized, try /createUser or /authorize", http.StatusBadRequest, user
	}

	return "you're authorized ", http.StatusOK, user
}
