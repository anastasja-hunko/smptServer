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

			h.serv.WriteResponse(rw, err.Error(), http.StatusBadRequest, user)

			return

		}
		cfg := mailjob.NewConfig()

		msg := "From: " + cfg.SmtpAddress + "\n" +
			"To: " + m.AddressTo + "\n" +
			"Subject: " + m.Header + "\n" +
			m.Body

		addr := cfg.SmtpServer + ":" + cfg.SmtpPort

		auth := smtp.PlainAuth("", cfg.SmtpAddress, cfg.SmtpPassword, cfg.SmtpServer)

		err = smtp.SendMail(addr, auth, cfg.SmtpAddress, []string{m.AddressTo}, []byte(msg))

		if err != nil {

			h.serv.WriteResponse(rw, err.Error(), http.StatusBadRequest, user)

			return
		}

		h.serv.DB.UserCol.UpdateUserMessages(user, m)

		h.serv.WriteResponse(rw, "message was sent", http.StatusOK, user)
	}
}
