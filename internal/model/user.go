package model

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

//User struct
type User struct {
	Login    string `json:"login, omitempty" bson:"_id"`
	Password string `json:"password, omitempty" bson:"password"`
	Active   bool   `json:"active" bson="active"`

	History  []History `json:"-" bson:"history"`
	UserLog  []Log     `json:"-" bson:"log"`
	Messages []Message `json:"-" bson:"messages"`
}

//HashPass encodes user's password
func (u *User) HashPass() error {

	if len(u.Password) > 0 {

		enc, err := hashPassword(u.Password)
		if err != nil {
			return err
		}

		u.Password = enc
	}

	return nil
}

//ComparePasswords returns equality entered password with user password
func (u *User) ComparePasswords(password string) bool {

	return checkPasswordHash(u.Password, password)
}

func checkPasswordHash(hash, password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

//hash password
func hashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

//AppendToHistoryAndLogs adds history and log records to user
func (u *User) AppendToHistoryAndLogs(field string, oldValue interface{}, newValue interface{}) (*[]History, *[]Log) {

	history := NewHistory(field, oldValue, newValue)

	histories := u.History

	if histories == nil {

		histories = []History{}

	}

	histories = append(histories, *history)

	logs := u.AppendToLogs(fmt.Sprintf("field %v was updated from %v to %v", field, oldValue, newValue))

	return &histories, logs
}

//AppendToLogs log record to user
func (u *User) AppendToLogs(text string) *[]Log {

	logs := u.UserLog

	if logs == nil {

		logs = []Log{}

	}

	message := NewLog(text)

	logs = append(logs, *message)

	return &logs
}

//AppendToMessages adds message record to user
func (u *User) AppendToMessages(msg *Message) *[]Message {

	messages := u.Messages

	if messages == nil {

		messages = []Message{}

	}

	messages = append(messages, *msg)

	return &messages
}
