package validation

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/fr"
	ut "github.com/go-playground/universal-translator"
	gpvalidator "github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	frTrans "github.com/go-playground/validator/v10/translations/fr"

	"go-nextjs-dashboard/internal/optional"
)

var (
	Validator *gpvalidator.Validate
	Uni       *ut.UniversalTranslator
)

func init() {
	Validator = gpvalidator.New(gpvalidator.WithRequiredStructEnabled())

	// registers a function to get alternate JSON names
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return fld.Name
		}
		return name
	})

	// setup translations
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

	registerOptionalType[string]()
	registerOptionalType[int]()
	registerOptionalType[bool]()
	registerOptionalType[float64]()

	Validator.RegisterAlias("rfc3339", "datetime="+time.RFC3339)
	err := Validator.RegisterTranslation(
		"rfc3339",
		enT,
		func(ut ut.Translator) error {
			return ut.Add("rfc3339", "{0} must be a valid RFC-3339 date-time", true)
		},
		func(ut ut.Translator, fe gpvalidator.FieldError) string {
			t, _ := ut.T("rfc3339", fe.Field())
			return t
		},
	)
	if err != nil {
		log.Fatal(fmt.Errorf("register translations error: %w", err))
	}
}

func registerOptionalType[T any]() {
	var z T
	t := reflect.TypeOf(z)

	// Check if T is a pointer
	if t.Kind() == reflect.Ptr {
		log.Fatalf("registerOptionalType: type parameter T must not be a pointer, got %v\n", t)
	}

	Validator.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		o, ok := field.Interface().(optional.Optional[T])

		if !ok {
			return nil
		}
		// not in JSON ⇒ return nil of T
		if !o.IsPresent {
			return new(T)
		}
		// explicitly null ⇒ treat as zero value
		if o.IsNull {
			var zero T
			return zero
		}

		// present + non-null ⇒ underlying value
		return o.Val
	}, optional.Optional[T]{})

	Validator.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		o, ok := field.Interface().(optional.Optional[*T])

		if !ok {
			return nil
		}
		// not in JSON or explicitly null ⇒ return current o.Val (which should be nil)
		if !o.IsPresent || o.IsNull {
			return o.Val
		}

		// present + non-null ⇒ underlying value
		return *o.Val
	}, optional.Optional[*T]{})
}
