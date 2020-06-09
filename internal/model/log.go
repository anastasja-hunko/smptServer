package model

import "time"

type Log struct {
	Text  string `bson:"text"`
	Login string `bson:"login"`
	Time  string `bson:"time"`
}

func NewLog(text string, login string) *Log {

	return &Log{
		Text:  text,
		Login: login,
	}
}

func (l *Log) BeforeCreate() {

	l.Time = time.Now().Format("Mon Jan _2 15:04:05 2006")
}
