package model

type History struct {
	Field    string      `json:"field" bson:"field"`
	OldValue interface{} `json:"oldValue" bson:"oldValue"`
	NewValue interface{} `json:"newValue" bson:"newValue"`
}

func NewHistory(field string, oldValue interface{}, newValue interface{}) *History {
	return &History{Field: field, OldValue: oldValue, NewValue: newValue}
}
