package models

type (
	Answer struct {
		Text      string `json:"text"`
		Letter    string `json:"letter"`
		IsCorrect bool   `json:"is_correct"`
	}
)
