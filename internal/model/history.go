package model

type History struct {
	field    string `json:"field" bson:"field"`
	oldValue string `json:"oldValue" bson:"oldValue"`
	newValue string `json:"newValue" bson:"newValue"`
}

func NewHistory(field string, oldValue string, newValue string) *History {
	return &History{field: field, oldValue: oldValue, newValue: newValue}
}
