package model

type Artwork struct {
	ID           int    `json:"objectID"`
	Title        string `json:"title"`
	Artist       string `json:"artistDisplayName"`
	Culture      string `json:"culture"`
	ObjectDate   string `json:"objectDate"`
	PrimaryImage string `json:"primaryImage"`
}
