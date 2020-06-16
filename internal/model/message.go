package model

type Message struct {
	AddressTo string `json:"to"`
	Header    string `json:"mailHeader"`
	Body      string `json:"mailBody"`
}
