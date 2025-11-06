package validation

import (
	"log"
	"reflect"

	erwp "github.com/PerumallaGiridhar/oolio/internal/errorwrap"
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

func ValidatePromocodeBloom(bf *bloom.BloomFilter) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		field := fl.Field()
		if field.Kind() != reflect.String {
			return false
		}
		code := field.String()
		if len(code) < 8 || len(code) > 10 {
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
		if len(code) < 8 || len(code) > 10 {
			return false
		}

		return idx.IsValid2of3(code)
	}

}

func RegisterPromocodeValidation(index *index.PebbleIndex) {
	log.Println("Registering bloom filter validator")
	erwp.MustDo(Validator.RegisterValidation("promocode", ValidatePromocodePebble(index)))
	log.Println("Initializing bloom filter completed")
}

func RegisterTranslations() {
	universalTranslator := ut.New(en.New())
	Translator, _ := universalTranslator.GetTranslator("en")
	erwp.MustDo(enTranslations.RegisterDefaultTranslations(Validator, Translator))
}

func HTTPRequestValidatorInit(index *index.PebbleIndex) {

	log.Println("Initializing request validator...")
	Validator = validator.New()

	RegisterTranslations()
	RegisterPromocodeValidation(index)

}
