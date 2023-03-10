package validators

import "reflect"

// NewIsTypeOf returns a function that returns true only if the interface type,
// determined by reflect.TypeOf(), is the same as the supplied exemplar.
// nil values always return false.
// The returned function can be used as a named func in the
// govalidators.CustomTypeTagMap for use in struct attribute value validations
func NewIsTypeOf(exemplar any) func(i interface{}, context interface{}) bool {
	t := reflect.TypeOf(exemplar).String()
	return func(i interface{}, context interface{}) bool {
		if i == nil {
			return false
		}
		return reflect.TypeOf(i).String() == t
	}
}
