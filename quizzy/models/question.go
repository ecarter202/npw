package models

import "encoding/json"

type (
	Question struct {
		Answers []*Answer `json:"answers"`

		ID   string `json:"id"`
		Text string `json:"text"`
	}
)

func (q *Question) JSON() string {
	b, err := json.MarshalIndent(q, "", "    ")
	if err != nil {
		return ""
	}

	return string(b)
}
