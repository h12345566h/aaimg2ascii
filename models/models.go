package models

import (
	"time"
)

type User struct {
	UserId    uint32 `gorm:"primary_key"`
	Account   string `gorm:"not null;type:varchar(20);unique"`
	Password  string `gorm:"not null;type:varchar(100)"`
	Email     string `gorm:"not null;type:varchar(100);unique"`
	Photo     string `gorm:"not null;type:varchar(100);"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type AsciiArt struct {
	AsciiArtId   uint32 `gorm:"primary_key"`
	AsciiContent string `gorm:"not null;type:text"`
	Public       bool   `gorm:"not null;type:boolean"`
	Row          int    `gorm:"not null;type:int"`
	Col          int    `gorm:"not null;type:int"`
	Hot          int    `gorm:"not null;type:int"`
	UserId      uint32
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
