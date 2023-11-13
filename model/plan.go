package model

import "time"

type Plan struct {
	Title         string         `json:"title" bson:"title"`
	Date          string         `json:"date" bson:"date"`
	CreatedAt     time.Time      `json:"created_at" bson:"created_at"`
	Absent        []Absent       `json:"absent" bson:"absent"`
	Substitutions []Substitution `json:"substitutions" bson:"substitutions"`
}

type Absent struct {
	Class   string `json:"class" bson:"class"`
	Periods string `json:"absent_time" bson:"periods"`
}

type Substitution struct {
	Class      string `json:"class" bson:"class"`
	Teacher    string `json:"teacher" bson:"teacher"`
	Period     string `json:"period" bson:"period"`
	Substitute string `json:"substitute,omitempty" bson:"substitute,omitempty"`
	Room       string `json:"room,omitempty" bson:"room,omitempty"`
	Info       string `json:"info,omitempty" bson:"info,omitempty"`
}
