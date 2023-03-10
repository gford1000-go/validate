package validators

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
)

func isFormattedAsUKPostcode(postcode string) bool {
	// Per: https://www.oreilly.com/library/view/regular-expressions-cookbook/9781449327453/ch04s16.html
	r := regexp.MustCompile("^[A-Z]{1,2}[0-9R][0-9A-Z]? [0-9][ABD-HJLNP-UW-Z]{2}$")
	return r.Match([]byte(strings.ToUpper(strings.TrimSpace(postcode))))
}

func isFormattedAsUSZipcode(zip string) bool {
	// Per: https://regexlib.com/Search.aspx?k=us+zip+code
	r := regexp.MustCompile("(^[0-9]{5}$)|(^[0-9]{9}$)|(^[0-9]{5}-[0-9]{4}$)")
	return r.Match([]byte(strings.ToUpper(strings.TrimSpace(zip))))
}

func TestPostcodeValidation(t *testing.T) {

	assert.True(t, isFormattedAsUKPostcode("AA9A 9AA"))
	assert.True(t, isFormattedAsUKPostcode("A9A 9AA"))
	assert.True(t, isFormattedAsUKPostcode("A9 9AA"))
	assert.True(t, isFormattedAsUKPostcode("A99 9AA"))
	assert.True(t, isFormattedAsUKPostcode("AA9 9AA"))
	assert.True(t, isFormattedAsUKPostcode("AA99 9AA"))

	assert.True(t, isFormattedAsUKPostcode("aa9A 9AA"))
	assert.True(t, isFormattedAsUKPostcode("a9A 9AA"))
	assert.True(t, isFormattedAsUKPostcode("a9 9AA"))
	assert.True(t, isFormattedAsUKPostcode("a99 9AA"))
	assert.True(t, isFormattedAsUKPostcode("aa9 9AA"))
	assert.True(t, isFormattedAsUKPostcode("aa99 9AA"))
	assert.True(t, isFormattedAsUKPostcode("aa9a 9AA"))
	assert.True(t, isFormattedAsUKPostcode("a9a 9AA"))
	assert.True(t, isFormattedAsUKPostcode("aa9a 9aa"))
	assert.True(t, isFormattedAsUKPostcode("a9A 9aa"))
	assert.True(t, isFormattedAsUKPostcode("a9 9aa"))
	assert.True(t, isFormattedAsUKPostcode("a99 9aa"))
	assert.True(t, isFormattedAsUKPostcode("aa9 9aa"))
}

func TestZipcodeValidation(t *testing.T) {

	assert.True(t, isFormattedAsUSZipcode("12345"))
	assert.True(t, isFormattedAsUSZipcode("123456789"))
	assert.True(t, isFormattedAsUSZipcode("12345-6789"))
	assert.False(t, isFormattedAsUSZipcode("12"))
	assert.False(t, isFormattedAsUSZipcode("123"))
	assert.False(t, isFormattedAsUSZipcode("1234"))
	assert.False(t, isFormattedAsUSZipcode("123456"))
	assert.False(t, isFormattedAsUSZipcode("1234567"))
	assert.False(t, isFormattedAsUSZipcode("12345678"))
}

// address is an example struct, where the Postcode validation needs to
// be performed conditionally on the Country value
type address struct {
	Country  string `json:"country" valid:"required,isUSUK"`
	Postcode string `json:"postcode" valid:"required,checkPostcode"`
}

// GetAsString required to implement Getter for conditional validation
func (a address) GetAsString(attrName string) (string, error) {
	switch attrName {
	case "Country":
		return a.Country, nil
	case "Postcode":
		return a.Postcode, nil
	default:
		return "", fmt.Errorf("%v unknown", attrName)
	}
}

func TestNewConditionalCheck(t *testing.T) {

	// This limits our tests to either UK or US as valid values of Country attribute
	govalidator.CustomTypeTagMap.Set("isUSUK", NewIsIn([]string{"UK", "US"}))

	// Wrap our postcode format testing code as a custom type validator
	validUK := func(i interface{}, context interface{}) bool {
		return isFormattedAsUKPostcode(fmt.Sprint(i))
	}
	govalidator.CustomTypeTagMap.Set("isUKPostcode", validUK)

	// Wrap our zipcode format testing code as a custom type validator
	validUS := func(i interface{}, context interface{}) bool {
		return isFormattedAsUSZipcode(fmt.Sprint(i))
	}
	govalidator.CustomTypeTagMap.Set("isUSZipcode", validUS)

	// Create the clauses to contextually validate the postcode details
	clauses := []ConditionClause{
		{
			ContextAttrName: "Country",
			ValueConditions: map[string][]string{
				"UK": {"isUKPostcode"},
				"US": {"isUSZipcode"},
			},
		},
	}

	// This registers the conditional validator
	govalidator.CustomTypeTagMap.Set("checkPostcode", NewConditionalCheck(clauses))

	type tc struct {
		m  string
		ok bool
	}

	tests := []tc{
		{
			`
{
	"country" : "UK",
	"postcode" : "EC4A 1HQ"
}`,
			true,
		},
		{
			`
{
	"country" : "US",
	"postcode" : "11123"
}`,
			true,
		},
		{
			`
{
	"country" : "UK",
	"postcode" : "11123"
}`,
			false,
		},
		{
			`
{
	"country" : "US",
	"postcode" : "M21 3GZ"
}`,
			false,
		},
		{
			`
{
	"country" : "SG",
	"postcode" : "M21 3GZ"
}`,
			false,
		},
		{
			`
{
	"postcode" : "M21 3GZ"
}`,
			false,
		},
		{
			`
{
	"country" : "UK"
}`,
			false,
		},
	}

	for _, test := range tests {
		var a = address{}
		json.Unmarshal([]byte(test.m), &a)

		_, err := govalidator.ValidateStruct(&a)

		if test.ok && err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !test.ok && err == nil {
			t.Fatal("Expected error, but passed")
		}
	}
}
