package internal

import (
	"github.com/anastasja-hunko/smptServer/internal/mailjob"
	"github.com/anastasja-hunko/smptServer/internal/model"
	"net/http"
	"net/smtp"
)

type sendHandler struct {
	serv *Server
}

func NewSendHandler(serv *Server) *sendHandler {

	return &sendHandler{serv: serv}

}

func (h *sendHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		h.send(rw, r)

		return

	}

	h.serv.Respond(rw, nil, "views/sendMail.html")

}

func (h *sendHandler) send(w http.ResponseWriter, r *http.Request) {

	m := &model.Message{
		AddressTo: r.FormValue("to"),
		Header:    r.FormValue("mailHeader"),
		Body:      r.FormValue("mailBody"),
	}

	cfg := mailjob.NewConfig()

	msg := "From: " + cfg.SmtpAddress + "\n" +
		"To: " + m.AddressTo + "\n" +
		"Subject: " + m.Header + "\n" +
		m.Body

	addr := cfg.SmtpServer + ":" + cfg.SmtpPort

	auth := smtp.PlainAuth("", cfg.SmtpAddress, cfg.SmtpPassword, cfg.SmtpServer)

	err := smtp.SendMail(addr, auth, cfg.SmtpAddress, []string{m.AddressTo}, []byte(msg))

	if err != nil {

		h.serv.writeErrorLog(err)

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	h.serv.writeOKMessage("Message was sent", "")

	http.Redirect(w, r, "/", 302)

	return

}
