package models

type (
	Answer struct {
		Text      string `json:"text"`
		IsCorrect bool   `json:"is_correct"`
		Letter    string `json:",omitempty"`
	}
)
