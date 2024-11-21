package testifyhelper

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/lumiluminousai/testify/mock"
)

type TestFunc func(t *testing.T) interface{}

func RunTest(testFunc TestFunc) func(t *testing.T) {
	return func(t *testing.T) {
		handler := testFunc(t)
		if handler != nil {
			if err := AssertExpectationsForMocks(t, handler); err != nil {
				t.Error(err)
			}
		}
	}
}

// MockTestingT is a custom implementation of TestingT that captures errors
type MockTestingT struct {
	Errors []string
}

func (m *MockTestingT) Errorf(format string, args ...interface{}) {
	m.Errors = append(m.Errors, fmt.Sprintf(format, args...))
}

func (m *MockTestingT) FailNow() {
	// No action needed; we are just collecting errors
}

func (m *MockTestingT) Logf(format string, args ...interface{}) {
	// No action needed; implement if needed
}

func AssertExpectationsForMocks(t *testing.T, handler interface{}) error {
	t.Helper()
	v := reflect.ValueOf(handler)

	// Ensure the handler is a pointer to a struct
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("AssertExpectationsForMocks requires a pointer to a struct")
	}

	// Dereference to access the struct
	v = v.Elem()

	return traverseFields(t, v)
}

func traverseFields(t *testing.T, v reflect.Value) error {
	t.Helper()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		// Skip unexported fields unless they are embedded
		if fieldType.PkgPath != "" && !fieldType.Anonymous {
			continue
		}

		fieldValue := field

		// If the field is a pointer, get its value
		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			fieldValue = fieldValue.Elem()
		}

		// Check if field is of type mock.Mock
		if fieldValue.Type() == reflect.TypeOf(mock.Mock{}) {
			mockField := fieldValue.Addr().Interface().(*mock.Mock)
			// Use MockTestingT to capture errors
			mt := &MockTestingT{}
			if !mockField.AssertExpectations(mt) {
				// Collect errors from MockTestingT
				errMessage := strings.Join(mt.Errors, "\n")
				// Simplify the error message
				err := fmt.Errorf("assert expectations failed for mock field '%s': %s", fieldType.Name, errMessage)
				return err
			}
			continue
		}

		// For structs, including embedded ones, recurse into their fields
		if fieldValue.Kind() == reflect.Struct {
			if err := traverseFields(t, fieldValue); err != nil {
				return err
			}
			continue
		}

		// Handle interfaces containing mocks
		if fieldValue.Kind() == reflect.Interface && !fieldValue.IsNil() {
			impl := fieldValue.Elem()
			if impl.Kind() == reflect.Ptr && impl.Elem().Kind() == reflect.Struct {
				if err := AssertExpectationsForMocks(t, impl.Interface()); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
