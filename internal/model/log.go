package model

import "time"

type Log struct {
	Text string `bson:"text"`
	Time string `bson:"time"`
}

func NewLog(text string) *Log {

	return &Log{
		Text: text,
	}
}

func (l *Log) BeforeCreate() {

	l.Time = time.Now().Format("Mon Jan _2 15:04:05 2006")
}
