package words_counters

type wordsCountersDto struct {
	Phrase  string `validate:"required,min=1,max=200" json:"phrase"`
	Counter int32  `validate:"gte=0" json:"counter"`
	Enabled *bool  `validate:"required" json:"enabled"`
}
