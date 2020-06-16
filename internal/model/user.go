package model

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Login    string `json:"login, omitempty" bson:"_id"`
	Password string `json:"password, omitempty" bson:"password"`
	Active   bool   `json:"active" bson="active"`

	History  []History `json:"-" bson:"history"`
	UserLog  []Log     `json:"-" bson:"log"`
	Messages []Message `json:"-" bson:"messages"`
}

//hash user's password and  before save to db
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

//compare password while autorizing
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

func (u *User) AppendToHistoryAndLogs(field string, oldValue interface{}, newValue interface{}) (*[]History, *[]Log) {

	history := NewHistory(field, oldValue, newValue)

	histories := []History{}

	if u.History != nil {

		histories = u.History

	}

	histories = append(histories, *history)

	logs := u.AppendToLogs(fmt.Sprintf("field %v was updated from %v to %v", field, oldValue, newValue))

	return &histories, logs
}

func (u *User) AppendToLogs(text string) *[]Log {
	logs := []Log{}

	if u.UserLog != nil {

		logs = u.UserLog

	}

	message := NewLog(text)

	logs = append(logs, *message)

	return &logs
}

func (u *User) AppendToMessages(msg *Message) *[]Message {
	messages := []Message{}

	if u.Messages != nil {

		messages = u.Messages

	}

	messages = append(messages, *msg)

	return &messages
}
