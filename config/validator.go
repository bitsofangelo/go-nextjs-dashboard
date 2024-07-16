package config

import (
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/fr"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"

	en_translations "github.com/go-playground/validator/v10/translations/en"
	fr_translations "github.com/go-playground/validator/v10/translations/fr"
)

// use a single instance of Validate, it caches struct info
var Validate *validator.Validate

var Uni *ut.UniversalTranslator

func InitValidator() {
	en := en.New()
	fr := fr.New()
	Uni = ut.New(en, en, fr)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see Uni.FindTranslator(...)
	trans, _ := Uni.GetTranslator("en")
	transFr, _ := Uni.GetTranslator("fr")

	Validate = validator.New(validator.WithRequiredStructEnabled())
	en_translations.RegisterDefaultTranslations(Validate, trans)
	fr_translations.RegisterDefaultTranslations(Validate, transFr)

	// register function to get tag name from json tags.
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	Validate.RegisterValidation("test100", func(fl validator.FieldLevel) bool {
		// User.Age needs to fit our needs, 12-18 years old.
		return fl.Field().Float() >= 100
	})

	// var regTransFunc validator.RegisterTranslationsFunc
	// var transFunc validator.TranslationFunc

	Validate.RegisterTranslation("test100", trans, registrationFunc("test100", "{0} must be above 100", false), translateFunc)

	// errs := Validate.VarWithKey("email", "test", "required,email")

	// if errs != nil {
	// 	log.Println(errs) // output: Key: "" Error:Field validation for "" failed on the "email" tag
	// 	return
	// }
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		if err = ut.Add(tag, translation, override); err != nil {
			return
		}

		return
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		log.Printf("warning: error translating FieldError: %#v", fe)
		return fe.(error).Error()
	}

	return t
}

func ValidateMapErr(data map[string]interface{}, rules map[string]interface{}) error {
	var errs validator.ValidationErrors
	mapErrs := Validate.ValidateMap(data, rules)

	if len(mapErrs) == 0 {
		return nil
	}

	for _, mapErr := range mapErrs {
		if fieldErrs, ok := mapErr.(validator.ValidationErrors); ok {
			errs = append(errs, fieldErrs...)
		}
	}

	return errs
}
