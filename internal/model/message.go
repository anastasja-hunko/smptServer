package model

//Message struct
type Message struct {
	AddressTo string `json:"to" bson:"to"`
	Header    string `json:"mailHeader" bson:"mailHeader"`
	Body      string `json:"mailBody" bson:"-"`
}
