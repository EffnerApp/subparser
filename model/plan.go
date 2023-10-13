package model

import "time"

type Plan struct {
	Title         string         `json:"title"`
	Date          string         `json:"date"`
	CreatedAt     time.Time      `json:"created_at"`
	Absent        []Absent       `json:"absent"`
	Substitutions []Substitution `json:"substitutions"`
}

type Absent struct {
	Class   string `json:"class"`
	Periods string `json:"absent_time"`
}

type Substitution struct {
	Class      string `json:"class"`
	Teacher    string `json:"teacher"`
	Period     string `json:"period"`
	Substitute string `json:"substitute,omitempty"`
	Room       string `json:"room,omitempty"`
	Info       string `json:"info,omitempty"`
}
