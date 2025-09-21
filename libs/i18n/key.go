package i18n

type TranslationKey[V any] interface {
	IsTranslationKey()
	GetPath() string
	GetPathSlice() []string
	GetVars() Vars
	SetVars(vars V) TranslationKey[V]
}
