package model

import "time"

//Log struct
type Log struct {
	Text string `bson:"text"`
	Time string `bson:"time"`
}

//NewLog returns initialized log
func NewLog(text string) *Log {

	return &Log{
		Text: text,
	}
}

//BeforeCreate adds time to log
func (l *Log) BeforeCreate() {

	l.Time = time.Now().Format("Mon Jan _2 15:04:05 2006")
}
