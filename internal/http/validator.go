package http

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/fr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	frTrans "github.com/go-playground/validator/v10/translations/fr"
)

var (
	Validator *validator.Validate
	Uni       *ut.UniversalTranslator
)

func init() {
	Validator = validator.New()

	// Registers a function to get alternate JSON names
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return fld.Name
		}
		return name
	})

	// Setup translations
	enLocale := en.New()
	frLocale := fr.New()
	Uni = ut.New(enLocale, enLocale, frLocale)

	enT, _ := Uni.GetTranslator("en")
	frT, _ := Uni.GetTranslator("fr")

	if err := enTrans.RegisterDefaultTranslations(Validator, enT); err != nil {
		log.Fatal(fmt.Errorf("register en translations error: %w", err))
	}

	if err := frTrans.RegisterDefaultTranslations(Validator, frT); err != nil {
		log.Fatal(fmt.Errorf("register fr translations error: %w", err))
	}
}
