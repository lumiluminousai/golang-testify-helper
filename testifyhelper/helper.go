package testifyhelper

import (
	"reflect"
	"testing"

	"github.com/lumiluminousai/testify/mock"
)

type TestFunc func(t *testing.T) interface{}

func RunTest(testFunc TestFunc) func(t *testing.T) {
	return func(t *testing.T) {
		handler := testFunc(t)
		if handler != nil {
			AssertExpectationsForMocks(t, handler)
		}
	}
}

func AssertExpectationsForMocks(t *testing.T, handler interface{}) {
	v := reflect.ValueOf(handler)

	// Ensure the handler is a pointer to a struct
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		t.Errorf("AssertExpectationsForMocks requires a pointer to a struct")
		return
	}

	// Dereference to access the struct
	v = v.Elem()
	t.Logf("Handler type: %s", v.Type())

	traverseFields(t, v)
}

func traverseFields(t *testing.T, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		t.Logf("Inspecting field: %s (%s)", fieldType.Name, field.Type())

		// Skip unexported fields unless they are embedded
		if fieldType.PkgPath != "" && !fieldType.Anonymous {
			t.Logf("Skipping unexported field: %s (%s)", fieldType.Name, field.Type())
			continue
		}

		fieldValue := field

		// If the field is a pointer, get its value
		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			fieldValue = fieldValue.Elem()
		}

		// Check if field is of type mock.Mock
		if fieldValue.Type() == reflect.TypeOf(mock.Mock{}) {
			t.Logf("Found mock field: %s (%s)", fieldType.Name, fieldValue.Type())
			mockField := fieldValue.Addr().Interface().(*mock.Mock)
			mockField.AssertExpectations(t)
			continue
		}

		// For structs, including embedded ones, recurse into their fields
		if fieldValue.Kind() == reflect.Struct {
			t.Logf("Recursing into struct field: %s (%s)", fieldType.Name, fieldValue.Type())
			traverseFields(t, fieldValue)
			continue
		}

		// Handle interfaces containing mocks
		if fieldValue.Kind() == reflect.Interface && !fieldValue.IsNil() {
			impl := fieldValue.Elem()
			if impl.Kind() == reflect.Ptr && impl.Elem().Kind() == reflect.Struct {
				t.Logf("Drilling down into interface implementation for: %s", fieldType.Name)
				AssertExpectationsForMocks(t, impl.Interface()) // Recurse for interface
			}
		}
	}
}
