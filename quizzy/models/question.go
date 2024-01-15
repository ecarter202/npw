package models

type (
	Question struct {
		Answers []*Answer `json:"answers"`

		Text string `json:"text"`
	}
)
