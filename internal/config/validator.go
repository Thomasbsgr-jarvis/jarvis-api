package config

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var (
	ErrFields = errors.New("Fields error:")
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
			case "max":
				return fmt.Errorf("%w Le champ %s ne peut pas dépasser %s caractères.", ErrFields, field.Field(), field.Param())
			case "url":
				return fmt.Errorf("%w L'URL est invalide.", ErrFields)
			case "uuid4":
				return fmt.Errorf("%w L'identifiant est invalide.", ErrFields)
			case "oneof":
				return fmt.Errorf("%w La valeur du champ %s est invalide.", ErrFields, field.Field())
			case "gte", "lte":
				return fmt.Errorf("%w La valeur du champ %s est hors limites.", ErrFields, field.Field())
			default:
				return fmt.Errorf("%w Champ invalide: %s", ErrFields, field.Field())
			}
		}
		return fmt.Errorf("erreur inattendue : %w", err)
	}
	return nil
}
