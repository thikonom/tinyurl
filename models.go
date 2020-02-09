package tinyurl

import (
        "time"
        // "github.com/jinzhu/gorm"
)
type User struct {
        Email string `gorm:"primary_key"`
        CreatedAt time.Time
        UpdatedAt time.Time
        Name string
        LastLogin *time.Time
        TinyURLS []TinyURL  `gorm:"ForeignKey:UserEmail"`
}

type TinyURL struct {
        ShortenedURL string `gorm:"primary_key"`
        CreatedAt time.Time
        OriginalURL string
        UserEmail string
}

