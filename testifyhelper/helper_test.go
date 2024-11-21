package testifyhelper

import (
	"testing"

	"github.com/lumiluminousai/testify/assert"
	"github.com/lumiluminousai/testify/mock"
	"github.com/lumiluminousai/testify/require"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) DoSomething(arg string) error {
	args := m.Called(arg)
	return args.Error(0)
}

type Handler struct {
	Service *MockService
}

func Test_AssertExpectationsForMocks_WithAssertNumberOfCalls_Success(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	mockService.On("DoSomething", "test").Return(nil).Once()

	handler := &Handler{
		Service: mockService,
	}

	// Act
	err := handler.Service.DoSomething("test")

	// Assert
	require.NoError(t, err)
	mockService.AssertNumberOfCalls(t, "DoSomething", 1)

	// This should pass without errors
	err = AssertExpectationsForMocks(t, handler)
	require.NoError(t, err)
}

func Test_AssertExpectationsForMocks_WithoutAssertNumberOfCalls_Failure(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	mockService.On("DoSomething", "test").Return(nil)

	handler := &Handler{
		Service: mockService,
	}

	// Act
	err := handler.Service.DoSomething("test")

	// Assert
	require.NoError(t, err)

	// Not calling AssertNumberOfCalls

	// This should fail and report missing AssertNumberOfCalls
	err = AssertExpectationsForMocks(t, handler)

	// Ensure err is not nil
	if assert.Error(t, err) {
		// Check the error message contains the expected message
		expectedErrorMessage := "assert expectations failed for mock field 'Mock': Missing AssertNumberOfCalls for method 'DoSomething'.\nFAIL: 0 out of 1 expectation(s) were met.\n\tThe code you are testing needs to make 1 more call(s)"
		assert.Contains(t, err.Error(), expectedErrorMessage)
	}
}

func Test_AssertExpectationsForMocks_MethodNotCalled_Failure(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	mockService.On("DoSomething", "test").Return(nil)

	handler := &Handler{
		Service: mockService,
	}

	// Act
	// Not calling handler.Service.DoSomething("test")

	// Assert
	// This should fail and report that the expected method was not called
	err := AssertExpectationsForMocks(t, handler)
	require.Error(t, err)
	require.Contains(t, err.Error(), "assert expectations failed")
}
