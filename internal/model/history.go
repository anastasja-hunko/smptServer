package model

type History struct {
	field    string      `json:"field" bson:"field"`
	oldValue interface{} `json:"oldValue" bson:"oldValue"`
	newValue interface{} `json:"newValue" bson:"newValue"`
}

func NewHistory(field string, oldValue interface{}, newValue interface{}) *History {
	return &History{field: field, oldValue: oldValue, newValue: newValue}
}
