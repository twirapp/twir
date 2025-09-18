package i18n

type TranslationKey interface {
	IsTranslationKey()
	GetPath() string
	GetPathSlice() []string
}
