package validation

import (
	"fmt"
	"log"
	"reflect"

	"github.com/PerumallaGiridhar/oolio/internal/index"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/willf/bloom"
)

var (
	Validator  *validator.Validate
	Translator ut.Translator
)

func isValidPromocodeLength(promocode string) bool {
	if len(promocode) < 8 || len(promocode) > 10 {
		return false
	}
	return true
}

func ValidatePromocodeBloom(bf *bloom.BloomFilter) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		field := fl.Field()
		if field.Kind() != reflect.String {
			return false
		}
		code := field.String()
		if !isValidPromocodeLength(code) {
			return false
		}

		return bf.TestString(code)
	}

}

func ValidatePromocodePebble(idx *index.PebbleIndex) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		field := fl.Field()
		if field.Kind() != reflect.String {
			return false
		}
		code := field.String()
		if !isValidPromocodeLength(code) {
			return false
		}

		validated, err := idx.IsValid2of3(code)
		if err != nil {
			log.Panicf("Error validating promocode")
		}

		return validated
	}

}

func RegisterPromocodeValidation(index *index.PebbleIndex) error {
	log.Println("Registering bloom filter validator")
	if err := Validator.RegisterValidation("promocode", ValidatePromocodePebble(index)); err != nil {
		return err
	}
	log.Println("Initializing bloom filter completed")

	return nil
}

func RegisterTranslations() error {
	universalTranslator := ut.New(en.New(), en.New())
	var ok bool
	Translator, ok = universalTranslator.GetTranslator("en")
	if !ok {
		return fmt.Errorf("no translator for 'en'")
	}
	if err := enTranslations.RegisterDefaultTranslations(Validator, Translator); err != nil {
		return err
	}

	return nil
}

func HTTPRequestValidatorInit(index *index.PebbleIndex) error {

	log.Println("Initializing request validator...")
	Validator = validator.New()

	if err := RegisterTranslations(); err != nil {
		return err
	}

	if err := RegisterPromocodeValidation(index); err != nil {
		return err
	}

	return nil
}
