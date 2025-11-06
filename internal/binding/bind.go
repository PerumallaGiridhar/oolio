package binding

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/PerumallaGiridhar/oolio/internal/validation"
	"github.com/go-playground/validator/v10"
)

func BindAndValidateJSONRequest(r *http.Request, dst any) map[string]string {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return map[string]string{"error": "invalid or unknown JSON fields"}
	}

	if err := validation.Validator.Struct(dst); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			fields := make(map[string]string, len(ve))
			log.Println("ve: ", ve)
			for _, fe := range ve {
				fields[fe.Field()] = fe.Translate(validation.Translator)
			}
			return fields
		}
		return map[string]string{"error": err.Error()}
	}

	return nil
}
