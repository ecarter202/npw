package models

type (
	Subject struct {
		Name string `json:"name"`
		// Version is for cert version (e.g. Network+ 008)
		Version string `json:"version"`
	}
)
