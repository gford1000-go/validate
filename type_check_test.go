package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIsTypeOf(t *testing.T) {

	var str = ""

	isString := NewIsTypeOf(str)
	isStringPtr := NewIsTypeOf(&str)
	isStringSlice := NewIsTypeOf([]string{})
	isMapStringInt := NewIsTypeOf(map[string]int{})

	var hello = "Hello"
	assert.True(t, isString(hello, nil))
	assert.False(t, isString(1, nil))

	// nil always fails
	assert.False(t, isString(nil, nil))

	// int(1) is not a string
	assert.False(t, isString(1, nil))

	// pointers are tested explicitly
	assert.False(t, isString(&hello, nil))
	assert.True(t, isStringPtr(&hello, nil))

	// slices
	assert.False(t, isStringSlice(str, nil))
	assert.True(t, isStringSlice([]string{str}, nil))

	// interesting
	m := map[string]int{"a": 1}
	var i interface{} = m
	assert.True(t, isMapStringInt(m, nil))
	assert.True(t, isMapStringInt(i, nil))

	type x struct{}

	isX := NewIsTypeOf(x{})
	assert.True(t, isX(x{}, nil))
	assert.False(t, isX(&x{}, nil))
}
