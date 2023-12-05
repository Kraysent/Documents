package validation

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
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

func In[T comparable](obj T, set []T) Rule {
	return func() error {
		if !slices.Contains(set, obj) {
			return NewValidationError(obj, fmt.Sprintf("value not in a set: %v", set))
		}

		return nil
	}
}

type Number interface {
	constraints.Integer | constraints.Float
}

func IsBetween[T Number](obj T, left T, right T) Rule {
	return func() error {
		if obj < left {
			return NewValidationError(obj, fmt.Sprintf("number is less than %d", left))
		}

		if obj > right {
			return NewValidationError(obj, fmt.Sprintf("number is bigger than %d", right))
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
