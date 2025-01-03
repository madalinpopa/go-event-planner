package validator

import (
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

// EmailRX is a compiled regular expression to validate the format of email addresses according to standard RFC 5322 rules.
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Validator is a struct used to validate data and store field-specific error messages.
type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Valid checks if the validator's FieldErrors map is empty, indicating that no validation errors are present.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// AddFieldError adds a validation error message for a specific field if it does not already exist in the FieldErrors map.
func (v *Validator) AddFieldError(key, message string) {

	// Check if the FieldErrors map is nil.
	// If it is nil, initialize it to a new empty map to avoid runtime errors when adding values.
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// AddNonFieldError appends a non-field-specific error message to the NonFieldErrors slice.
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// CheckField adds a validation error for the specified
// field if the condition is not met.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank checks if the provided string is not empty or whitespace-only
// and returns true if it contains non-whitespace characters.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars checks if the number of characters in a string is
// less than or equal to a specified limit.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedValue checks if a given value is present within a
// list of permitted values and returns true if found.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// ValidDate checks if the provided date string represents a valid,
// non-zero date and returns true if valid, otherwise false.
func ValidDate(date time.Time) bool {
	if date.IsZero() {
		return false
	}
	return true
}

// MinChars checks if the string `value` contains at least `n` characters,
// returning true if the condition is met.
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Matches checks if the given string value matches
// the specified regular expression pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
