package validator

import (
	"bytes"
	"go-nextjs-dashboard/config"
	"strings"

	"github.com/gofiber/fiber/v2/log"

	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"
)

type MapValidationErrors map[string][]error

func (e MapValidationErrors) Error() string {
	buff := bytes.NewBufferString("")

	for key, errs := range e {
		for _, err := range errs {
			buff.WriteString(key + ": " + err.Error())
			buff.WriteString("\n")
		}
	}

	return strings.TrimSpace(buff.String())
}

// Translate translates all of the ValidationErrors
func (e MapValidationErrors) Translate(ut ut.Translator) map[string][]string {
	trans := make(map[string][]string, len(e))

	if len(e) == 0 {
		return trans
	}

	for key, mapErrs := range e {
		for _, err := range mapErrs {
			if fieldErr, ok := err.(validator.FieldError); ok {
				t, err := ut.T(fieldErr.Tag(), key, fieldErr.Param())
				if err != nil {
					log.Errorf("warning: error translating FieldError: %#v", fieldErr)
					trans[key] = append(trans[key], fieldErr.(error).Error())
				}
				trans[key] = append(trans[key], t)
			}
		}

		// if fieldErr, ok := mapErr.(validator.FieldError); ok {
		// 	// ut.T(key, key)
		// 	// trans[key] = fieldErr.Translate(ut)
		// 	t, err := ut.T(fieldErr.Tag(), key, fieldErr.Param())
		// 	if err != nil {
		// 		log.Errorf("warning: error translating FieldError: %#v", fieldErr)
		// 		trans[key] = fieldErr.(error).Error()
		// 	}
		// 	trans[key] = t
		// }
	}

	return trans
}

func ValidateMap(data map[string]interface{}, rules map[string]interface{}) error {
	// var errs validator.ValidationErrors
	mapErrs := config.Validate.ValidateMap(data, rules)

	if len(mapErrs) == 0 {
		return nil
	}

	newMapErrs := make(MapValidationErrors, len(mapErrs))

	for key, mapErr := range mapErrs {
		// newMapErrs[key] = make([]error)
		if fieldErrs, ok := mapErr.(validator.ValidationErrors); ok {
			for _, fieldErr := range fieldErrs {
				newMapErrs[key] = append(newMapErrs[key], fieldErr)
			}
		}
	}

	return newMapErrs
}

func ValidateMapQueries(data map[string]string, rules map[string]interface{}) error {
	queryData := make(map[string]any, len(data))

	for key, value := range data {
		queryData[key] = value
	}

	return ValidateMap(queryData, rules)
}
