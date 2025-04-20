package model

import "sync"

type Artwork struct {
	ID           int     `json:"objectID"`
	Title        *string `json:"title"`
	Artist       *string `json:"artistDisplayName"`
	Culture      *string `json:"culture"`
	ObjectDate   *string `json:"objectDate"`
	PrimaryImage *string `json:"primaryImage"`
}

type Cache struct {
	Mu       sync.RWMutex
	CacheMap map[int]*Artwork
}
