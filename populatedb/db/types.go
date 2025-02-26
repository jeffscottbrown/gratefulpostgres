package db

import "time"

type SongPerformance struct {
	ID         uint `gorm:"primarykey"`
	SetID      uint
	Title      string
	OrderInSet int
}

type Set struct {
	ID               uint `gorm:"primarykey"`
	ShowID           uint
	SetNumber        int
	SongPerformances []SongPerformance `gorm:"foreignKey:SetID"`
}

type Show struct {
	ID    uint      `gorm:"primarykey"`
	Date  time.Time `gorm:"type:date"`
	Venue string
	City  string
	State string
	Sets  []Set `gorm:"foreignKey:ShowID"`
}

type YamlShow struct {
	Venue   string           `yaml:":venue"`
	City    string           `yaml:":city"`
	State   string           `yaml:":state"`
	Country string           `yaml:":country"`
	Setlist []map[string]any `yaml:":sets"`
}
