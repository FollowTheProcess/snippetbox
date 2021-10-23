package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// Form represents a HTML form
type Form struct {
	url.Values
	Errors errors
}

// New initialises and returns a new Form
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required is the error raised when a required form field is blank
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MaxLength will add to the error pool if a field is over it's maximum
// char length
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

// PermittedValues will add to the error pool if a field has an invalid choice
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// Valid returns true if the form is valid, else false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
