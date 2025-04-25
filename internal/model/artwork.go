package model

import "sync"

type Artwork struct {
	ID           int     `json:"objectID"`
	LastUsedAt   int64   //最終利用時刻
	Title        *string `json:"title"`
	Artist       *string `json:"artistDisplayName"`
	Culture      *string `json:"culture"`
	ObjectDate   *string `json:"objectDate"`
	PrimaryImage *string `json:"primaryImage"`
}

type Cache struct {
	Mu       sync.RWMutex
	CacheMap map[int]*Artwork
	MaxSize  int
	Ttl      int
}
