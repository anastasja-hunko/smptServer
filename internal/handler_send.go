package internal

import (
	"encoding/json"
	"github.com/anastasja-hunko/smptServer/internal/mailjob"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"net/http"
	"net/smtp"
)

type sendHandler struct {
	serv *Server
}

func newSendHandler(serv *Server) *sendHandler {

	return &sendHandler{serv: serv}

}

func (h *sendHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	_, _, user := h.serv.getInfoForRespond(r)

	if user != nil {

		var m = &model.Message{}

		err := json.NewDecoder(r.Body).Decode(m)
		if err != nil {

			h.serv.writeResponse(rw, err.Error(), http.StatusBadRequest, user)

			return

		}
		cfg := mailjob.NewConfig()

		msg := "From: " + cfg.SMTPAddress + "\n" +
			"To: " + m.AddressTo + "\n" +
			"Subject: " + m.Header + "\n" +
			m.Body

		addr := cfg.SMTPServer + ":" + cfg.SMTPPort

		auth := smtp.PlainAuth("", cfg.SMTPAddress, cfg.SMTPPassword, cfg.SMTPServer)

		err = smtp.SendMail(addr, auth, cfg.SMTPAddress, []string{m.AddressTo}, []byte(msg))
		if err != nil {

			h.serv.writeResponse(rw, err.Error(), http.StatusBadRequest, user)

			return
		}

		h.serv.DB.UserCol.UpdateUserMessages(user, m)

		h.serv.writeResponse(rw, "message was sent", http.StatusOK, user)
	}
}
