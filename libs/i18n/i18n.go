package i18n

import (
	"fmt"
	"strings"

	kaptinlini18n "github.com/kaptinlin/go-i18n"
)

// LocalesStore example: {"en": {"commands": {"followage": {"title": "qwe"}}}}
type LocalesStore map[string]map[string]map[string]map[string]string

type Opts struct {
	Store         LocalesStore
	DefaultLocale string
}

func New(opts Opts) (*I18n, error) {
	i := &I18n{
		store: opts.Store,
	}

	if opts.DefaultLocale == "" {
		return nil, fmt.Errorf("default locale is required")
	}

	if _, ok := opts.Store[opts.DefaultLocale]; !ok {
		return nil, fmt.Errorf("default locale %s not found in store", opts.DefaultLocale)
	}

	supportedLocales := make([]string, 0, len(opts.Store))
	for locale := range opts.Store {
		supportedLocales = append(supportedLocales, locale)
	}

	bundle := kaptinlini18n.NewBundle(
		kaptinlini18n.WithDefaultLocale(opts.DefaultLocale),
		kaptinlini18n.WithLocales(supportedLocales...),
	)
	i.bundle = bundle

	preparedBundle := map[string]map[string]string{}
	for locale, translations := range opts.Store {
		if _, ok := preparedBundle[locale]; !ok {
			preparedBundle[locale] = map[string]string{}
		}

		flatTranslations := map[string]string{}
		for subMapKey, subMap := range translations {
			for subSubMapKey, subSubMap := range subMap {
				for key, value := range subSubMap {
					flatTranslations[fmt.Sprintf("%s_%s_%s", subMapKey, subSubMapKey, key)] = value
				}
			}
		}
		preparedBundle[locale] = flatTranslations
	}

	if err := bundle.LoadMessages(preparedBundle); err != nil {
		return nil, fmt.Errorf("failed to load messages: %v", err)
	}

	return i, nil
}

type I18n struct {
	store         LocalesStore
	defaultLocale string
	bundle        *kaptinlini18n.I18n
}

type Vars map[string]interface{}

func (i *I18n) T(locale string, key TranslationKey, vars ...Vars) string {
	localizer := i.bundle.NewLocalizer(locale)

	convertedVars := make(kaptinlini18n.Vars)
	for _, v := range vars {
		for key, value := range v {
			convertedVars[key] = value
		}
	}

	return localizer.Get(strings.Join(key.GetPathSlice(), "_"), convertedVars)
}
