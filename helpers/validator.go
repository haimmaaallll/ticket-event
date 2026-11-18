package helpers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(data any) map[string]string {
	//cek struct sesuai tag `validate:"..."`
	err := validate.Struct(data)
	if err == nil {
		// semua isian valid
		return nil
	}

	// ubah error umum menjadi error validasi per field
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// bukan error validasi (kasus langka), kembalikan apa adanya
		return map[string]string{"error": err.Error()}
	}

	// kumpulkan pesan error, kuncinya nama field
	result := make(map[string]string)
	for _, fieldError := range validationErrors {
		switch fieldError.Tag() {
		case "required":
			result[fieldError.Field()] = "wajib diisi"
		case "email":
			result[fieldError.Field()] = "format email tidak valid"
		case "min":
			result[fieldError.Field()] = fmt.Sprintf("minimal %s karakter", fieldError.Param())
		case "max":
			result[fieldError.Field()] = fmt.Sprintf("maksimal %s karakter", fieldError.Param())
		default:
			result[fieldError.Field()] = "isian tidak valid"
		}
	}

	return result
}
