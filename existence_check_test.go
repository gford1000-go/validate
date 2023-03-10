package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateExistenceCheck(t *testing.T) {

	equals1 := NewEquals(1)
	in0to5 := NewIsIn([]int{0, 1, 2, 3, 4, 5})

	var one = 1
	var six = 6

	assert.True(t, equals1(one, nil))    // 1 == 1
	assert.True(t, equals1(&one, nil))   // 1 == 1
	assert.False(t, equals1(six, nil))   // 6 != 1
	assert.False(t, equals1("one", nil)) // "one" != 1
	assert.False(t, equals1(nil, nil))   // nil != 1
	assert.True(t, in0to5(one, nil))     // 1 in [0..5]
	assert.True(t, in0to5(&one, nil))    // 1 in [0..5]
	assert.False(t, in0to5(six, nil))    // 6 not in [0..5]
	assert.True(t, in0to5([]int{1, 2, 3}, nil))
	assert.True(t, in0to5(&[]int{1, 2, 3}, nil))
	assert.False(t, in0to5([]int{-1, 2, 3}, nil))
	assert.False(t, in0to5(&[]int{-1, 2, 3}, nil))
	assert.True(t, in0to5([]int{1, 2, 3, 4, 5}, nil))
	assert.False(t, in0to5([]int{1, 2, 3, 4, 6}, nil))

	notEquals1 := NewNotEquals(1)
	notIn0to5 := NewIsNotIn([]int{0, 1, 2, 3, 4, 5})

	assert.False(t, notEquals1(one, nil)) // 1 == 1
	assert.True(t, notEquals1(2, nil))    // 2 != 1
	assert.True(t, notIn0to5(6, nil))
	assert.True(t, notIn0to5([]int{6, 8, 9}, nil))

	testVals := []int{1, 6, 8, 9}
	testPtrVals := []*int{&testVals[0], &testVals[1], &testVals[2], &testVals[3], nil}

	// nil never matches
	assert.False(t, equals1(testPtrVals[4], nil))
	assert.False(t, notEquals1(testPtrVals[4], nil))

	// equality test applied through pointer values
	assert.True(t, equals1(testPtrVals[0], nil))
	assert.False(t, in0to5(testPtrVals[1], nil))

	assert.False(t, notEquals1(testPtrVals[0], nil))
	assert.True(t, notIn0to5(testPtrVals[1], nil))

	// nil never matches within a slice as well
	assert.False(t, in0to5([]*int{testPtrVals[0], testPtrVals[4]}, nil))
	assert.False(t, notIn0to5([]*int{testPtrVals[1], testPtrVals[4]}, nil))

	assert.True(t, in0to5([]*int{testPtrVals[0], testPtrVals[0]}, nil))
	assert.True(t, notIn0to5([]*int{testPtrVals[1], testPtrVals[2]}, nil))

	assert.True(t, in0to5(&[]*int{testPtrVals[0], testPtrVals[0]}, nil))
	assert.True(t, notIn0to5(&[]*int{testPtrVals[1], testPtrVals[2]}, nil))
}
