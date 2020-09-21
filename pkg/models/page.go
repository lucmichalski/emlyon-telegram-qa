package models

import (
	"fmt"
	"time"

	"github.com/abadojack/whatlanggo"
	"gorm.io/gorm"
)

type Page struct {
	gorm.Model
	LinkHash           string    `json:"hash,omitempty"`
	Title              string    `json:"title,omitempty"`
	Body               string    `json:"body,omitempty"`
	MetaDescription    string    `json:"description,omitempty"`
	MetaKeywords       string    `json:"keywords,omitempty"`
	MetaLang           string    `json:"lang,omitempty"`
	CanonicalLink      string    `json:"canonical-link,omitempty"`
	CleanedText        string    `json:"cleaned-text,omitempty"`
	FinalURL           string    `json:"url,omitempty"`
	TopImage           string    `json:"image,omitempty"`
	PublishDate        time.Time `json:"publish-date,omitempty"`
	Language           string    `json:"detected-lang,omitempty"`
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
