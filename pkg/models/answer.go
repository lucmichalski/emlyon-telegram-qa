package models

type HaystackParams struct {
	Question      string `url:"question"`
	Index         string `url:"index"`
	TopKReader    int    `url:"top_k_reader"`
	TopKRetriever int    `url:"top_k_retriever"`
}

type Haystack struct {
	Answers  []HaystackAnswer `json:"answers"`
	NoAnsGap float64          `json:"no_ans_gap"`
	Question string           `json:"question"`
}

type HaystackAnswer struct {
	Answer           string             `json:"answer"`
	Context          string             `json:"context"`
	DocumentID       string             `json:"document_id"`
	Meta             HaystackAnswerMeta `json:"meta"`
	OffsetEnd        int                `json:"offset_end"`
	OffsetEndInDoc   int                `json:"offset_end_in_doc"`
	OffsetStart      int                `json:"offset_start"`
	OffsetStartInDoc int                `json:"offset_start_in_doc"`
	Probability      float64            `json:"probability"`
	Score            float64            `json:"score"`
}

type HaystackAnswerMeta struct {
	CanonicalLink   string   `json:"canonical-link"`
	Image           string   `json:"image"`
	Language        string   `json:"language"`
	MetaDescription string   `json:"meta-description"`
	MetaKeywords    []string `json:"meta-keywords"`
	MetaLang        string   `json:"meta-lang"`
	Name            string   `json:"name"`
	PublishDate     string   `json:"publish-date"`
	Summary         string   `json:"summary"`
	URL             string   `json:"url"`
	Yake            []string `json:"yake"`
}
