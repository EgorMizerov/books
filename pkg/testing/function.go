package testing

import (
	"reflect"

	"github.com/stretchr/testify/assert"
)

// FuncsEqual checks if expectedFunc and actualFunc are pointing the same function,
// it will return false on equivalent, but not the same functions.
func FuncsEqual(t assert.TestingT, expectedFunc, actualFunc interface{}) {
	expected := reflect.ValueOf(&expectedFunc).Elem()
	actual := reflect.ValueOf(&actualFunc).Elem()
	assert.True(t, expected.InterfaceData() == actual.InterfaceData())
}
