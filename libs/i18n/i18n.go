package i18n

import (
	"context"
	"fmt"
	"strings"

	kaptinlini18n "github.com/kaptinlin/go-i18n"
)

type Opts struct {
	Store         LocalesStore
	DefaultLocale string
}

// this is global because it is just much easier to use with GetCtx generic function
var globali18nInstance *I18n

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
					// replace dots from flattened nested keys with underscores to match path generation
					normalizedKey := strings.ReplaceAll(key, ".", "_")
					flatTranslations[fmt.Sprintf("%s_%s_%s", subMapKey, subSubMapKey, normalizedKey)] = value
				}
			}
		}
		preparedBundle[locale] = flatTranslations
	}

	if err := bundle.LoadMessages(preparedBundle); err != nil {
		return nil, fmt.Errorf("failed to load messages: %v", err)
	}

	globali18nInstance = i

	return i, nil
}

type I18n struct {
	store         LocalesStore
	defaultLocale string
	bundle        *kaptinlini18n.I18n
}

type Vars map[string]any

type LocaleCtxKey struct{}

type translateOptions struct {
	locale string
}

type TranslateOption = func(*translateOptions)

func WithLocale(locale string) TranslateOption {
	return func(options *translateOptions) {
		options.locale = locale
	}
}

func SetContextLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, LocaleCtxKey{}, locale)
}

func GetCtx[T any](
	ctx context.Context,
	key TranslationKey[T],
	opts ...TranslateOption,
) string {
	locale, ok := ctx.Value(LocaleCtxKey{}).(string)
	if ok && locale != "" {
		opts = append(opts, WithLocale(locale))
	}

	return Get(key, opts...)
}

func Get[T any](
	key TranslationKey[T],
	opts ...TranslateOption,
) string {
	if globali18nInstance == nil {
		panic("i18n instance is not initialized")
	}

	o := translateOptions{
		locale: globali18nInstance.defaultLocale,
	}

	for _, option := range opts {
		option(&o)
	}

	keyPath := strings.Join(key.GetPathSlice(), "_")
	localizer := globali18nInstance.bundle.NewLocalizer(o.locale)
	convertedVars := make(kaptinlini18n.Vars)
	for key, value := range key.GetVars() {
		convertedVars[key] = fmt.Sprint(value)
	}

	return localizer.Get(keyPath, convertedVars)
}
