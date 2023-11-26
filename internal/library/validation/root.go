package validation

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Rule func() error

func StringNotEmpty(str string) Rule {
	return func() error {
		if str == "" {
			return NewValidationError(str, "string is not empty")
		}

		return nil
	}
}

func IsUUID(str string) Rule {
	return func() error {
		_, err := uuid.Parse(str)
		if err != nil {
			return NewValidationError(str, "string is not a valid uuid")
		}

		return nil
	}
}

func IsISO8601(str string) Rule {
	return func() error {
		_, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return NewValidationError(str, "string is not a valid ISO 8601 timestamp")
		}

		return nil
	}
}

type Rules struct {
	rules []Rule
}

func NewRules(rules ...Rule) *Rules {
	return &Rules{
		rules: rules,
	}
}

func (r *Rules) Validate() error {
	errs := make([]ValidationError, 0)

	for _, rule := range r.rules {
		err := rule()
		var targetErr ValidationError

		if errors.As(err, &targetErr) {
			errs = append(errs, targetErr)
			continue
		}

		if err != nil {
			return fmt.Errorf("error during validation: %s", err.Error())
		}
	}

	if err := mergeErrors(errs); !err.IsEmpty() {
		return err
	}

	return nil
}
