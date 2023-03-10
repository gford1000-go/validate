package validators

import (
	"fmt"
	"reflect"

	"github.com/asaskevich/govalidator"
)

// Getter provides a mechanism to access the value of particular attributes, by name
type Getter interface {
	GetAsString(attrName string) (string, error)
}

// ConditionClause defines which custom validators should be executed, based on the
// value of the specified attribute within the current context.
type ConditionClause struct {
	ContextAttrName string              // The name of the attribute in the context, whose value to retrieve
	ValueConditions map[string][]string // The set of validations to apply based on the attribute's value
}

// NewConditionalCheck creates a custom type validator which applies an arbitray number of
// ConditionClause checks, each with an arbitrary number of validators.  The context object
// must support the Getter interface to apply the ConditionClause validations ... if it doesn't
// then the custom type validator returns false.
//
//The validators within the ConditionClause slices must have been pre-registered to
// govalidator CustomTypeTagMap; if not found the custom type validator will panic.
func NewConditionalCheck(clauses []ConditionClause) func(i interface{}, context interface{}) bool {

	// Iterate through the validators to confirm i is validated in the context
	validate := func(i interface{}, context interface{}, validators []string) bool {

		passed := true
		for _, validatorName := range validators {
			validator, ok := govalidator.CustomTypeTagMap.Get(validatorName)
			if !ok {
				// This is unexpected and would lead to unpredictable behaviour - so panic
				panic(fmt.Sprintf("Validator %s is not in CustomTypeTagMap - cannot continue", validatorName))
			}

			passed = validator(i, context)
			if !passed {
				break
			}
		}
		return passed
	}

	// Iterate through the clauses, identifying the appropriate validator slice to
	// apply for each clause, based on the value of the clause's attribute
	eval := func(o Getter, i interface{}, context interface{}) bool {

		passed := true
		for _, clause := range clauses {

			v, err := o.GetAsString(clause.ContextAttrName)
			if err != nil {
				// This is unexpected and would lead to unpredictable behaviour - so panic
				panic(fmt.Sprintf("Attribute %s does not exist in %s, cannot validate (err: %v)", clause.ContextAttrName, reflect.TypeOf(o), err))
			}

			validators, ok := clause.ValueConditions[v]
			if ok {
				passed = validate(i, context, validators)
			}

			if !passed {
				break
			}
		}

		return passed
	}

	return func(i interface{}, context interface{}) bool {
		if i == nil || context == nil {
			return false
		}

		switch o := context.(type) {
		case Getter:
			return eval(o, i, context)
		default:
			fmt.Println(reflect.TypeOf(o))
			return false
		}
	}
}
