package config

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/fr"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"

	entranslations "github.com/go-playground/validator/v10/translations/en"
	frtranslations "github.com/go-playground/validator/v10/translations/fr"
)

// Validate use a single instance of Validate, it caches struct info
var Validate *validator.Validate

var Uni *ut.UniversalTranslator

func InitValidator() {
	enLocale := en.New()
	frLocale := fr.New()
	Uni = ut.New(enLocale, enLocale, frLocale)

	// this is usually know or extracted from storage 'Accept-Language' header
	// also see Uni.FindTranslator(...)
	trans, _ := Uni.GetTranslator("en")
	transFr, _ := Uni.GetTranslator("fr")

	Validate = validator.New(validator.WithRequiredStructEnabled())

	if err := entranslations.RegisterDefaultTranslations(Validate, trans); err != nil {
		log.Fatal(fmt.Errorf("register en translations error: %w", err))
	}

	if err := frtranslations.RegisterDefaultTranslations(Validate, transFr); err != nil {
		log.Fatal(fmt.Errorf("register fr translations error: %w", err))
	}

	// register function to get tag name from json tags.
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	_ = Validate.RegisterValidation("test100", func(fl validator.FieldLevel) bool {
		// User.Age needs to fit our needs
		return fl.Field().Float() >= 100
	})

	_ = Validate.RegisterTranslation("test100", trans, registrationFunc("test100", "{0} must be above 100", false), translateFunc)

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
		return fe.Error()
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
