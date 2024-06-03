package validator

import "regexp"

type ValidatorToy struct {
	Errors map[string]string
}

func New() *ValidatorToy {
	return &ValidatorToy{Errors: make(map[string]string)}
}

func (v *ValidatorToy) Valid() bool {
	return len(v.Errors) == 0
}

func (v *ValidatorToy) AddErrorToy(key, message string) {
	if _, exist := v.Errors[key]; !exist {
		v.Errors[key] = message
	}
}

func (v *ValidatorToy) Check(ok bool, key, message string) {
	if !ok {
		v.AddErrorToy(key, message)
	}
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}
