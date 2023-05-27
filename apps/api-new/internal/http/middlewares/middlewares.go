package middlewares

import (
	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"reflect"
	"strings"
)

type Middlewares struct {
	logger     *zap.SugaredLogger
	validator  *validator.Validate
	translator ut.Translator
}

func newValidator() (*validator.Validate, ut.Translator) {
	v := validator.New()
	en := en_US.New()
	uni := ut.New(en, en)
	transEN, _ := uni.GetTranslator("en_US")
	err := enTranslations.RegisterDefaultTranslations(v, transEN)
	if err != nil {
		zap.S().Fatalln(err)
	}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return v, transEN
}

type Opts struct {
	fx.In

	Logger *zap.SugaredLogger
}

func NewMiddlewares(logger *zap.SugaredLogger) *Middlewares {
	v, translator := newValidator()

	return &Middlewares{
		logger:     logger,
		validator:  v,
		translator: translator,
	}
}
