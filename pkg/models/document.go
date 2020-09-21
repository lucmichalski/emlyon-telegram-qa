package models

import (
        "gorm.io/gorm"
        // "github.com/abadojack/whatlanggo"
)

type Document struct {
   gorm.Model
   Title string
   Content string
   Format string
}

