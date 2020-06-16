package model

import "golang.org/x/crypto/bcrypt"

//type UserLog

type User struct {
	Login    string `json:"login" bson:"_id"`
	Password string `json:"password" bson:"password"`
	Active   bool   `json:"active" bson="active"`
	History  *[]History
	//UserLog []*UserLog `json:"logs bson="logs"`
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

func CreateHistory(field string, oldValue interface{}, newValue interface{}) History {
	return NewHistory(field, oldValue, newValue)
}

//func CreateLog() {
//
//}

func (u *User) AppendToHistory(field string, oldValue interface{}, newValue interface{}) []History {
	history := CreateHistory(field, oldValue, newValue)

	histories := []History{}

	if u.History != nil {

		histories = *u.History

	}

	return append(histories, history)
}
