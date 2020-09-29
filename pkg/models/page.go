package models

import (
	"fmt"
	"time"

	"github.com/abadojack/whatlanggo"
	"gorm.io/gorm"
)

type Page struct {
	gorm.Model
	LinkHash           string    `json:"hash"`
	Title              string    `json:"title"`
	Body               string    `json:"body"`
	MetaDescription    string    `json:"description"`
	MetaKeywords       string    `json:"keywords"`
	MetaLang           string    `json:"lang"`
	CanonicalLink      string    `json:"canonical-link"`
	CleanedText        string    `json:"cleaned-text"`
	FinalURL           string    `json:"url"`
	TopImage           string    `json:"image"`
	PublishDate        time.Time `json:"publish-date"`
	Language           string    `json:"detected-lang"`
	Fingerprint        string    `json:"fingerprint"`
	LanguageConfidence float64   `json:"-"`
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
