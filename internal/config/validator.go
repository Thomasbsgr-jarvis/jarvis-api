package config

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateData(input any) error {
	if err := validate.Struct(input); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			field := validationErrors[0]
			switch field.Tag() {
			case "required":
				return fmt.Errorf("%w Tous les champs sont requis.", ErrFields)
			case "email":
				return fmt.Errorf("%w L'adresse email est invalide.", ErrFields)
			case "min":
				return fmt.Errorf("%w Le champ %s doit contenir au moins %s caractères.", ErrFields, field.Field(), field.Param())
			default:
				return fmt.Errorf("%w Champ invalide: %s", ErrFields, field.Field())
			}
		}
		return fmt.Errorf("erreur inattendue : %w", err)
	}
	return nil
}
