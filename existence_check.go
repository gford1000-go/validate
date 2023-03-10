package validators

// ExistenceCheckType describes the style of the existence check
type existenceCheckType int

const (
	ExhaustiveIn    existenceCheckType = iota // All items are found in the comparison slice
	ExhaustiveNotIn                           // None of the items are in the comparison slice
)

func (e existenceCheckType) String() string {
	v := []string{"ExhaustiveIn", "ExhaustiveNotIn"}
	return v[e]
}

// createExistenceCheck returns an instance of a validation function that can be added to the
// govalidators.CustomTypeTagMap for use in struct attribute value validations
func createExistenceCheck[T comparable](checkType existenceCheckType, compareWith []T) func(i interface{}, context interface{}) bool {

	compSlice := func(v []T) bool {
		statusOfCompare := false
		outcomeForError := false
		for _, item := range v {
			switch checkType {
			case ExhaustiveIn:
				statusOfCompare = false
			case ExhaustiveNotIn:
				statusOfCompare = true
			default:
				panic("Invalid checkType")
			}

			for _, ci := range compareWith {
				if item == ci {
					statusOfCompare = !statusOfCompare
					break
				}
			}
			if statusOfCompare == outcomeForError {
				return false
			}
		}
		return true
	}

	compPtrSlice := func(v []*T) bool {
		for _, item := range v {
			if item == nil {
				return false
			}
			if !compSlice([]T{*item}) {
				return false
			}
		}
		return true
	}

	return func(i interface{}, context interface{}) bool {
		if i == nil {
			return false
		}

		switch v := i.(type) {
		case T:
			return compSlice([]T{v})
		case []T:
			return compSlice(v)
		case *[]T:
			return compSlice(*v)
		case *T:
			return compPtrSlice([]*T{v})
		case []*T:
			return compPtrSlice(v)
		case *[]*T:
			return compPtrSlice(*v)
		default:
			return false
		}
	}
}

// NewEquals returns a function instance that will test equality of an interface against the specified item.
// The returned function can be used as a named func in the
// govalidators.CustomTypeTagMap for use in struct attribute value validations
func NewEquals[T comparable](item T) func(i interface{}, context interface{}) bool {
	return createExistenceCheck(ExhaustiveIn, []T{item})
}

// NewNotEquals returns a function instance that will test non-equality of an interface against the specified item.
// The returned function can be used as a named func in the
// govalidators.CustomTypeTagMap for use in struct attribute value validations
func NewNotEquals[T comparable](item T) func(i interface{}, context interface{}) bool {
	return createExistenceCheck(ExhaustiveNotIn, []T{item})
}

// NewIsIn returns a function instance that will test existence of an interface against the specified items.
// The function will return true if the interface is a value and is present, or when the interface is a slice
// and all the values in the slice are present in items.
// The returned function can be used as a named func in the
// govalidators.CustomTypeTagMap for use in struct attribute value validations
func NewIsIn[T comparable](items []T) func(i interface{}, context interface{}) bool {
	return createExistenceCheck(ExhaustiveIn, items)
}

// NewIsNotIn returns a function instance that will test non-existence of an interface against the specified items.
// The function will return true if the interface is a value and is not present, or when the interface is a slice
// and none of the values in the slice are present in items.
// The returned function can be used as a named func in the
// govalidators.CustomTypeTagMap for use in struct attribute value validations
func NewIsNotIn[T comparable](items []T) func(i interface{}, context interface{}) bool {
	return createExistenceCheck(ExhaustiveNotIn, items)
}
