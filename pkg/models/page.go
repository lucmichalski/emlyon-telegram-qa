package models

import (
	"fmt"

	"gorm.io/gorm"
	"github.com/abadojack/whatlanggo"
)

type Page struct {
   gorm.Model
   Title string
   Body string 
   Language string
LanguageConfidence float64
}

func (p *Page) BeforeCreate() (err error) {
	// add to whatlango
	if p.Body != "" {
		info := whatlanggo.Detect(p.Body)
		p.Language = info.Lang.String()
		p.LanguageConfidence = info.Confidence
		fmt.Println("======> Language:", p.Language, " Script:", whatlanggo.Scripts[info.Script], " Confidence: ", p.LanguageConfidence)
	}
	return
}

func (p *Page) AfterCreate() (err error) {
	// add to manticore
	// add to bleve
	return
}
