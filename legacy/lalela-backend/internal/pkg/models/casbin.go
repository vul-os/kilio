package models

import "time"

//Role struct declaration
type Role struct {
	ID   uint   `gorm:"AUTO_INCREMENT"`
	Name string `json:"Name"`
}

//Group struct declaration
type Group struct {
	ID   uint   `gorm:"AUTO_INCREMENT"`
	Name string `json:"Name"`
}

//Action struct declaration
type Action struct {
	ID   uint   `gorm:"AUTO_INCREMENT"`
	Name string `json:"Name"`
}

//Event struct declaration
type Event struct {
	ID     uint      `gorm:"AUTO_INCREMENT"`
	UserID string    `json:"subject"`
	Event  string    `json:"event"`
	Object string    `json:"object"`
	Type   string    `json:"type"`
	Date   time.Time `json:"date"`
}

//Cashbin struct declaration
type CasbinRule struct {
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

// SiteConfig
type SiteConfig struct {
	ID    uint   `gorm:"AUTO_INCREMENT"`
	Rule  string `json:"subject"`
	Value bool
}
