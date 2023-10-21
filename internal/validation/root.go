package validation

import (
	"errors"
	"fmt"
)

type Rule func() error

func StringNotEmpty(str string) Rule {
	return func() error {
		if str == "" {
			return NewValidationError("string '%s' is empty", str)
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
